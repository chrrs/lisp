package main

import (
	"bufio"
	"fmt"
	"lisp/lisp"
	"os"
)

func main() {
	fmt.Println("clisp REPL (Ctrl-C to exit)")

	scanner := bufio.NewScanner(os.Stdin)

	env := lisp.NewEnvironment(nil)
	env.AddBuiltins()

	std, err := os.ReadFile("lib/std.clsp")
	if err == nil {
		_, err := Evaluate(&env, string(std), true)
		if err != nil {
			fmt.Printf("error while evaluating std.clsp: %v", err)
			return
		}
	}

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		out, err := Evaluate(&env, input, false)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(out)
		}
	}
}

func Evaluate(env *lisp.Environment, input string, multi bool) (lisp.Node, error) {
	tokens, err := lisp.Tokenize(input)
	if err != nil {
		return nil, fmt.Errorf("tokenization error: %v", err)
	}

	if multi {
		expression, err := lisp.ParseExpression(tokens, lisp.QExpression)
		if err != nil {
			return nil, fmt.Errorf("parsing error: %v", err)
		}

		for _, node := range expression.Nodes {
			out := node.Evaluate(env)
			err, ok := out.(lisp.ErrorNode)
			if ok {
				return nil, err.Error
			}
		}

		return lisp.ExpressionNode{Type: lisp.SExpression, Nodes: make([]lisp.Node, 0)}, nil
	} else {
		expression, err := lisp.ParseExpression(tokens, lisp.SExpression)
		if err != nil {
			return nil, fmt.Errorf("parsing error: %v", err)
		}

		return expression.Evaluate(env), nil
	}
}
