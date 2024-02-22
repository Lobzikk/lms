package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type Opers struct {
	Sum  int
	Div  int
	Prod int
	Sub  int
}

type ExpressionString struct {
	Exp string
}

type MainServer struct {
	Opers map[string]int
	Agent AgentServer
}

func (ms *MainServer) Start() {
	var (
		restartChan     chan struct{}
		shutdownChan    chan struct{}
		expressionsChan chan string
	)
	ms.Agent = *NewAgentServer()
	go ms.Agent.Start(restartChan, shutdownChan, expressionsChan, ms.Opers, false)
	mux := http.NewServeMux()
	mux.HandleFunc("/opers", func(w http.ResponseWriter, r *http.Request) {
		opersMap := make(map[string]int, 4)
		opersMap["+"], _ = strconv.Atoi(r.URL.Query().Get("Sum"))
		opersMap["-"], _ = strconv.Atoi(r.URL.Query().Get("Sub"))
		opersMap["*"], _ = strconv.Atoi(r.URL.Query().Get("Prod"))
		opersMap["/"], _ = strconv.Atoi(r.URL.Query().Get("Div"))
		ms.Opers = opersMap
		ms.Agent.mu.Lock()
		for i := 0; i < len(ms.Agent.MathServers); i++ {
			ms.Agent.MathServers[i].Opers = ms.Opers
		}
		ms.Agent.mu.Unlock()
	})
	mux.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		shutdownChan <- struct{}{}
	})
	mux.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
		restartChan <- struct{}{}
	})
	mux.HandleFunc("/newExpression", func(w http.ResponseWriter, r *http.Request) {
		exp := r.URL.Query().Get("Exp")
		ok, err := regexp.MatchString("[-+*/0-9()]", exp)
		if err != nil || !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Print(exp)
		expressionsChan <- exp
	})
	mux.HandleFunc("/getExpressions", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(ms.Agent.MathServers)
		if err != nil {
			panic(err)
		}
		w.Write(data)
	})
	http.ListenAndServe(":8000", mux)
}
