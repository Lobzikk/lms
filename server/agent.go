package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	epxressions "lms/expressions"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AgentServer struct {
	connection           *mongo.Client
	MathServers          []MathServer
	ServerGoroutineLimit int
	ServerAmount         int
	mu                   *sync.Mutex
	ErrorsChan           chan error
}

func NewAgentServer() *AgentServer {
	getEnv := func(key string, defaultValue int) int {
		value, err := strconv.Atoi(os.Getenv(key))
		if err != nil && value < 1 {
			return defaultValue
		}
		return value
	}
	amount := getEnv("COUNTING_SERVERS", 3)
	limit := getEnv("COUNTING_SERVERS_MAX_GOROUTINES", 10)
	return &AgentServer{
		ServerGoroutineLimit: limit,
		ServerAmount:         amount,
		mu:                   &sync.Mutex{},
	}
}

func GenNames(path string, amount int) []string {
	//read line by line names from names.txt
	names := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(names) == amount {
			break
		}
		names = append(names, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(names) < amount {
		origlen := len(names)
		for i := 0; i < amount-origlen; i++ {
			id := int(math.Mod(float64(i), float64(origlen)))
			names = append(names, names[id]+"_"+strconv.Itoa(i))
		}
	}
	return names
}

func (a *AgentServer) Fill(opers map[string]int, noDB bool) {
	a.MathServers = make([]MathServer, 0)
	if noDB {
		names := GenNames("server/names.txt", a.ServerAmount)
		for _, name := range names {
			a.MathServers = append(a.MathServers, *NewServer(name, opers))
		}
	} else {
		//Importing servers from Mongo DB
		if a.connection == nil {
			a.MongoConnect()
		}
		dbName := os.Getenv("MONGODB")
		db := a.connection.Database(dbName)
		serversCollection := db.Collection("servers")
		cur, err := serversCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			panic(err)
		}
		var res []interface{}
		if err := cur.All(context.TODO(), &res); err != nil {
			panic(err)
		}
		arr := make([]MathServer, len(res))
		for ind, server := range res {
			_, data, err := bson.MarshalValue(server)
			if err != nil {
				panic(err)
			}
			if err := bson.Unmarshal(data, &arr[ind]); err != nil {
				panic(err)
			}
		}
		a.MathServers = arr
	}
}

func (a *AgentServer) MongoConnect() {
	var err error
	uri := os.Getenv("MONGO_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	a.connection, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
}

func (a *AgentServer) ExportMathServers() { //(and clear the previous data)
	if a.connection == nil {
		a.MongoConnect()
	}
	dbName := os.Getenv("MONGODB")
	fmt.Println("a" + dbName)
	db := a.connection.Database(dbName)
	serversCollection := db.Collection("servers")
	serversCollection.DeleteMany(context.TODO(), nil)
	jsons := make([][]byte, len(a.MathServers))
	a.mu.Lock()
	for ind, server := range a.MathServers {
		json, err := json.Marshal(server)
		if err != nil {
			panic(err)
		}
		jsons[ind] = json
	}
	a.mu.Unlock()
	b := make([]interface{}, len(jsons))
	for ind, json := range jsons {
		err := bson.UnmarshalExtJSON(json, true, &b[ind])
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(b)
	_, err := serversCollection.InsertMany(context.TODO(), b)
	if err != nil {
		panic(err)
	}
}

func (a *AgentServer) Start(restart, shutdown chan struct{}, expressions chan string, opers map[string]int, restarted bool) {
	if opers == nil {
		opers = make(map[string]int, 4)
		var err error
		opers["SUM"], err = strconv.Atoi(os.Getenv("SUM"))
		if err != nil {
			panic(err)
		}
		opers["PROD"], err = strconv.Atoi(os.Getenv("PROD"))
		if err != nil {
			panic(err)
		}
		opers["SUB"], err = strconv.Atoi(os.Getenv("MINUS"))
		if err != nil {
			panic(err)
		}
		opers["DIV"], err = strconv.Atoi(os.Getenv("DIV"))
		if err != nil {
			panic(err)
		}
	}
	a.Fill(opers, !restarted)
	var id int
	expChanells := make([]chan epxressions.MathExpression, len(a.MathServers))
	killChannels := make([]chan struct{}, len(a.MathServers))
	for ind, server := range a.MathServers {
		go server.Start(expChanells[ind], killChannels[ind])
	}
	if restarted {
		a.mu.Lock()
		for _, server := range a.MathServers {
			for _, expression := range server.Expressions {
				if expression.ID > id {
					id = expression.ID //finding the biggest ID
				}
			}
		}
		a.mu.Unlock()
	}
	for {
		select {
		case <-restart:
			for _, channel := range killChannels {
				channel <- struct{}{}
			}
			fmt.Println("Server restarting in 15 seconds!")
			time.Sleep(15 * time.Second)
			a.Start(restart, shutdown, expressions, opers, true)
		case <-shutdown:
			for _, channel := range killChannels {
				channel <- struct{}{}
			}
			fmt.Println("Server shutting down in 15 seconds!")
			time.Sleep(15 * time.Second)
			a.ExportMathServers()
			return
		case expression := <-expressions:
			a.mu.Lock()
			id++
			sort.Slice(a.MathServers, func(i, j int) bool {
				return len(a.MathServers[i].Expressions) < len(a.MathServers[j].Expressions)
			})
			var full bool = true
			for ind, server := range a.MathServers {
				if len(server.Expressions) != a.ServerGoroutineLimit {
					full = false
					expChanells[ind] <- epxressions.MathExpression{
						Expression:  expression,
						SolvingTime: 0,
						ID:          id,
						Code:        200,
						IsMarked:    false,
					}
					break
				}
			}
			if full {
				a.ErrorsChan <- fmt.Errorf("error: server is currently full")
			}
			a.mu.Unlock()
		case <-time.After(10 * time.Second): //autosave every 10 seconds
			a.mu.Lock()
			a.ExportMathServers()
			os.Setenv("SUM", strconv.Itoa(opers["+"]))
			os.Setenv("DIV", strconv.Itoa(opers["/"]))
			os.Setenv("PROD", strconv.Itoa(opers["*"]))
			os.Setenv("SUB", strconv.Itoa(opers["-"]))
			a.mu.Unlock()
		}
	}
}
