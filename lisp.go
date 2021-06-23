package main

import (
	"bufio"
	"fmt"
	"lisp/lisp"
	"os"
)

func main() {
	fmt.Println("clisp v0")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		tokens, err := lisp.Tokenize(input)
		if err != nil {
			fmt.Println("tokenization error:", err)
			continue
		}

		expression, err := lisp.ParseExpression(tokens)
		if err != nil {
			fmt.Println("parsing error:", err)
			continue
		}

		expression.Dump(0)
	}
}
