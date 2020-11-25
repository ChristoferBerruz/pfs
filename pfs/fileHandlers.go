package pfs

import (
	"fmt"
	"os"
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
