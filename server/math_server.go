package server

import (
	expressions "lms/expressions"
	"slices"
	"sync"
	"time"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

type MathServer struct {
	name        string
	expressions []expressions.MathExpression
	opers       map[string]int
	error       bool
}

func NewServer(name string, opers map[string]int) *MathServer {
	return &MathServer{
		name:        name,
		expressions: make([]expressions.MathExpression, 0),
		error:       false,
		opers:       opers,
	}
}

func (ms *MathServer) Start(ch chan expressions.MathExpression) {
	var mu sync.Mutex
	for {
		select {
		case expression := <-ch:
			ms.expressions = append(ms.expressions, expression)
		case <-time.After(1 * time.Second):
			for ind, expression := range ms.expressions {
				if expression.Code == 400 {
					if !expression.IsMarked {
						go func() {
							expression.IsMarked = true
							//erase that element
							time.Sleep(10 * time.Second)
							mu.Lock()
							ind := slices.Index(ms.expressions, expression)
							ms.expressions[ind] = ms.expressions[len(ms.expressions)-1]
							ms.expressions[len(ms.expressions)-1] = expressions.MathExpression{}
							ms.expressions = ms.expressions[:len(ms.expressions)-1]
							mu.Unlock()
						}()
					}
					continue
				}
				ms.expressions[ind].SolvingTime++
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
					for oper := range ms.opers {
						if found {
							break
						}
						for ind, token := range postfix {
							if ind > 1 && token.Value.(string) == oper && postfix[ind-1].Type == 1 && postfix[ind-2].Type == 1 {
								expression.Current = postfix[ind-2].Value.(string) + token.Value.(string) + postfix[ind-1].Value.(string)
							} else if found {
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
