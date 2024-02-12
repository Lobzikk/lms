package math

import (
	"errors"
	"fmt"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

type MathExpression struct {
	Expression  string `json:"expression"`
	Current     string `json:"current"`
	SolvingTime int    `json:"solving_time"` // sec
	ID          int    `json:"id"`
	Code        int    `json:"code"` // http code
	IsMarked    bool   // whether the expression is marked or not (will be deleted in <10 sec)
}

func evaluateOperator(oper string, a, b int) (int, error) {
	switch oper {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, errors.New("unknown operator: " + oper)
	}
}

func Evaluate(tokens []*shuntingYard.RPNToken) (int, error) {
	if tokens == nil {
		return 0, errors.New("tokens cannot be nil")
	}
	var stack []int
	for _, token := range tokens {
		fmt.Println(stack, token.Value)
		// push all operands to the stack
		if token.Type == shuntingYard.RPNTokenTypeOperand {
			val := token.Value.(int)
			stack = append(stack, val)
		} else {
			// execute current operator
			if len(stack) < 2 {
				return 0, errors.New("missing operand")
			}
			// pop 2 elements
			arg1, arg2 := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			val, err := evaluateOperator(token.Value.(string), arg1, arg2)
			if err != nil {
				return 0, err
			}
			// push result back to stack
			stack = append(stack, val)
		}
	}
	if len(stack) != 1 {
		return 0, errors.New("stack corrupted")
	}
	return stack[len(stack)-1], nil
}
