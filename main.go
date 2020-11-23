package main

import (
	"fmt"
	"os"

	fcb "github.com/ChristoferBerruz/portable_file_system/fcb"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func createFile(filename string) {
	fd, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	if err := fd.Truncate(1e4); err != nil {
		panic(err)
	}

	fmt.Println("File successfully created.")
}

func main() {

	filename := ".pfs"
	if !fileExists(filename) {
		fmt.Println("PFS not found. Creating PFS...")
		createFile(filename)
	}

	block := fcb.NewFCB("test.txt", 32, 0, 0)
	fmt.Println(block)
	fmt.Println("Portable file system!")
}
