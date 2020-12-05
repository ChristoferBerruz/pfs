package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ChristoferBerruz/portable_file_system/pfs"
)

var fileSystem = pfs.FileSystem{}

func handleCommands(args []string) {
	switch args[0] {
	case "open":
		(&fileSystem).OpenVolume(args[1])
	case "put":
		(&fileSystem).Put(args[1])
	case "get":
		fmt.Println("get")
	case "rm":
		(&fileSystem).RemoveFile(args[1])
	case "dir":
		(&fileSystem).Dir()
	case "putr":
		fmt.Println("Putr...")
	case "kill":
		(&fileSystem).Kill(args[1])
	default:
		fmt.Println("Command is not valid!")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	shellPrompt := "pfs > "
	for {
		fmt.Print(shellPrompt)
		shellInput, _ := reader.ReadString('\n')
		shellInput = strings.TrimRight(shellInput, "\r\n")
		args := strings.Split(shellInput, " ")
		if args[0] == "quit" {
			if fileSystem.PfsFile != nil {
				(&fileSystem).Quit()
			}
			break
		}
		handleCommands(args)
		//pfs.TestDirectory()
	}
}
