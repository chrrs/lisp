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
