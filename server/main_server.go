package server

import (
	"encoding/json"
	"net/http"
	"regexp"
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
	ms.Agent.Start(restartChan, shutdownChan, expressionsChan, ms.Opers, false)
	mux := http.NewServeMux()
	mux.HandleFunc("/opers", func(w http.ResponseWriter, r *http.Request) {
		var opers Opers
		err := json.NewDecoder(r.Body).Decode(&opers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		opersMap := make(map[string]int, 4)
		opersMap["+"] = opers.Sum
		opersMap["-"] = opers.Sub
		opersMap["*"] = opers.Prod
		opersMap["/"] = opers.Div
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
		var exp ExpressionString
		if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, err := regexp.MatchString("[-+*/0-9()]", exp.Exp)
		if err != nil || !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		expressionsChan <- exp.Exp
	})
	mux.HandleFunc("/getExpressions", func(w http.ResponseWriter, r *http.Request) {
		ms.Agent.mu.Lock()
		data, err := json.Marshal(ms.Agent.MathServers)
		if err != nil {
			panic(err)
		}
		w.Write(data)
		ms.Agent.mu.Unlock()
	})
	go http.ListenAndServe(":8000", mux)
}
