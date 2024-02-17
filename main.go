package main

import (
	"fmt"
	expressions "lms/expressions"
	server "lms/server"
	"regexp"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

// executes an operator

func main() {
	input := "(2 /1)+(2*2) +   (2+2)*2"
	regex := "[-+*/0-9()]"
	fmt.Println(regexp.MatchString(regex, input))
	// parse input expression to infix notation
	infixTokens, err := shuntingYard.Scan(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Infix Tokens:")
	fmt.Println(infixTokens)

	// convert infix notation to postfix notation(RPN)
	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		panic(err)
	}
	fmt.Println("Postfix(RPN) Tokens:")
	for _, t := range postfixTokens {
		fmt.Printf("%v ", t.Type)
	}
	fmt.Println()

	// evaluate RPN tokens
	result, err := expressions.Evaluate(postfixTokens)
	if err != nil {
		panic(err)
	}

	// output the result
	fmt.Printf("Result: %v\n", result)
	fmt.Println(server.GenNames("server/names.txt", 8))
}
