package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ChristoferBerruz/portable_file_system/pfs"
)

var fileSystem = &pfs.FileSystem{}

func argsHasNArguments(args []string, n int) bool {
	return len(args) == n
}

func handleCommands(args []string) {
	switch args[0] {
	case "open":
		if argsHasNArguments(args, 2) {
			fileSystem.OpenVolume(args[1])
		}
	case "put":
		if argsHasNArguments(args, 2) {
			fileSystem.Put(args[1])
		}
	case "get":
		if argsHasNArguments(args, 2) {
			fileSystem.MoveOut(args[1])
		}
	case "rm":
		if argsHasNArguments(args, 2) {
			fileSystem.RemoveFile(args[1])
		}
	case "dir":
		if argsHasNArguments(args, 1) {
			fileSystem.Dir()
		}
	case "putr":
		if argsHasNArguments(args, 3) {
			fileSystem.PutRemarks(args[1], args[2])
		}
	case "kill":
		if argsHasNArguments(args, 2) {
			fileSystem.Kill(args[1])
		}
	default:
		fmt.Println("Command is not valid!")
	}
}

func keyboardInterrupHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fileSystem.Quit()
		os.Exit(0)
	}()
}

func main() {

	// In case we close abruptly
	keyboardInterrupHandler()

	reader := bufio.NewReader(os.Stdin)
	shellPrompt := "pfs > "
	for {
		fmt.Print(shellPrompt)
		shellInput, _ := reader.ReadString('\n')
		shellInput = strings.TrimRight(shellInput, "\r\n")
		args := strings.Split(shellInput, " ")

		// In case of empty enter
		if args[0] == "" {
			continue
		}

		if args[0] == "quit" {
			if (*fileSystem).PfsFile != nil {
				fileSystem.Quit()
			}
			break
		}
		handleCommands(args)
		//pfs.TestDirectory()
	}
}
