package server

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AgentServer struct {
	connection           *mongo.Client
	MathServers          []MathServer
	ServerGoroutineLimit int
	ServerAmount         int
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
		fmt.Println("aa")
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
	if noDB {
		names := GenNames("server/names.txt", a.ServerAmount)
		for _, name := range names {
			a.MathServers = append(a.MathServers, *NewServer(name, opers))
		}
	} else {
		//TODO: NoSQL connection to database with prepared statements
	}
}

func (a *AgentServer) MongoConnect() {
	uri := os.Getenv("MONGO_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(nil)
}
