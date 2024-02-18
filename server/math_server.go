package server

import (
	expressions "lms/expressions"
	"slices"
	"sync"
	"time"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

type MathServer struct {
	Name        string
	Expressions []expressions.MathExpression
	Opers       map[string]int
}

func NewServer(Name string, opers map[string]int) *MathServer {
	return &MathServer{
		Name:        Name,
		Expressions: make([]expressions.MathExpression, 0),
		Opers:       opers,
	}
}

func (ms *MathServer) Start(ch chan expressions.MathExpression, cancel chan struct{}) {
	var mu sync.Mutex
	for {
		select {
		case <-cancel:
			time.Sleep(11 * time.Second)
			return
		case expression := <-ch:
			ms.Expressions = append(ms.Expressions, expression)
		case <-time.After(1 * time.Second):
			for ind, expression := range ms.Expressions {
				if expression.SolvingTime > ms.Opers[string(expression.Current[1])] { //expression expired
					expression.Code = 503
				}
				if expression.Code == 400 || expression.Code == 503 {
					if !expression.IsMarked {
						go func() {
							expression.IsMarked = true
							//erase that element
							time.Sleep(10 * time.Second)
							mu.Lock()
							ind := slices.Index(ms.Expressions, expression)
							ms.Expressions[ind] = ms.Expressions[len(ms.Expressions)-1]
							ms.Expressions[len(ms.Expressions)-1] = expressions.MathExpression{}
							ms.Expressions = ms.Expressions[:len(ms.Expressions)-1]
							mu.Unlock()
						}()
					}
					continue
				}
				ms.Expressions[ind].SolvingTime++
				// operations hash map should be sorted in ascending order by this moment
				if expression.Current == "" {
					infix, err := shuntingYard.Scan(expression.Expression)
					if err != nil {
						expression.Code = 400
						break
					}
					postfix, err := shuntingYard.Parse(infix)
					if err != nil {
						expression.Code = 400
						break
					}
					var found bool
					for oper := range ms.Opers {
						if found {
							break
						}
						for ind, token := range postfix {
							if ind > 1 && token.Value.(string) == oper && postfix[ind-1].Type == 1 && postfix[ind-2].Type == 1 {
								expression.Current = postfix[ind-2].Value.(string) + token.Value.(string) + postfix[ind-1].Value.(string)
								found = true
							}
							if found {
								break
							}
						}
					}
				}
				//checking if expression is solvable
				if string(expression.Current[1]) == "/" && string(expression.Current[2]) == "0" {
					expression.Code = 400
				}
			}
		}
	}
}
