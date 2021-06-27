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
		_, err := lisp.Evaluate(&env, string(std), true)
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
		out, err := lisp.Evaluate(&env, input, false)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(out)
		}
	}
}
