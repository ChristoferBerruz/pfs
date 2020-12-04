package pfs

import (
	"fmt"
	"os"
)

// FileExists checks whethere the fileName exists in current directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// CreateFile creates a 10KB volume
func CreateFile(filename string) error {
	fd, err := os.Create(filename)

	if err != nil {
		return err
	}

	if err := fd.Truncate(1e4); err != nil {
		return err
	}

	fmt.Printf("Sucessfully created %s volume\n", filename)
	return nil
}
