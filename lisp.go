package main

import (
	"bufio"
	"fmt"
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
		fmt.Println(input)
	}
}
