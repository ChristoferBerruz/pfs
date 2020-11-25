package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		shellInput, _ := reader.ReadString('\n')
		shellInput = strings.TrimRight(shellInput, "\r\n")
		args := strings.Split(shellInput, " ")
		if args[0] == "exit" {
			break
		}

		fmt.Printf("Executing command ... %s\n", args[0])

	}
}
