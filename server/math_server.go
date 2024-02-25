package server

import (
	"fmt"
	expressions "lms/expressions"
	//"slices"
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
	for {
		select {
		case <-cancel:
			time.Sleep(11 * time.Second)
			return
		case expression := <-ch:
			ms.Expressions = append(ms.Expressions, expression)
			fmt.Print(ms)
		case <-time.After(1 * time.Second):
			for indExp, expression := range ms.Expressions {
				if expression.Current != "" {
					//checking if expression is solvable
					if string(expression.Current[1]) == "/" && string(expression.Current[2]) == "0" {
						expression.Code = 400
					}
					if expression.SolvingTime > ms.Opers[string(expression.Current[1])] { //expression expired
						expression.Code = 503
					}
					if expression.Code == 400 || expression.Code == 503 {
						// if !expression.IsMarked {
						// 	go func() {
						// 		expression.IsMarked = true
						// 		//erase that element
						// 		time.Sleep(10 * time.Second)
						// 		ind := slices.Index(ms.Expressions, expression)
						// 		ms.Expressions[ind] = ms.Expressions[len(ms.Expressions)-1]
						// 		ms.Expressions[len(ms.Expressions)-1] = expressions.MathExpression{}
						// 		ms.Expressions = ms.Expressions[:len(ms.Expressions)-1]
						// 	}()
						// }
						continue
					}
					ms.Expressions[indExp].SolvingTime++
					//checking if expression is solvable
					if string(expression.Current[1]) == "/" && string(expression.Current[2]) == "0" {
						expression.Code = 400
					}
				} else {
					// operations hash map should be sorted in ascending order by this moment
					infix, err := shuntingYard.Scan(expression.Expression)
					if err != nil {
						expression.Code = 400
						break
					}
					fmt.Print(infix)
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
							if ind > 1 && token.Value == oper && postfix[ind-1].Type == 1 && postfix[ind-2].Type == 1 {
								ms.Expressions[indExp].Current = fmt.Sprintf("%d%s%d", postfix[ind-2].Value, token.Value, postfix[ind-1].Value)
								found = true
								fmt.Println(expression.Current)
							}
							if found {
								break
							}
						}
					}
				}
			}
		}
	}
}
