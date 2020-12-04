package pfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// Directory contains data structures and metadata for the volume
type Directory struct {
	Metadata           Metadata
	FreeDataBlockArray [totalDataBlocks]bool
	FCBArray           [totalDataBlocks]FCB
}

// NewDirectory creates pfs directory
func NewDirectory(volName string) Directory {
	metadata := NewMetadata(volName)
	directory := Directory{
		Metadata: metadata,
	}
	for idx := range directory.FreeDataBlockArray {
		directory.FreeDataBlockArray[idx] = true
	}

	return directory
}

// ReadDirectoryFromDisk reads a directory from Disk
func ReadDirectoryFromDisk(file *os.File) (Directory, error) {

	// Directory is always from start of volume file
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Dumping data from file to buffer
	buf := make([]byte, directorySize)
	file.Read(buf)

	// Reading from buffer to create a directory
	var directory Directory
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, &directory); err != nil {
		return Directory{}, err
	}
	return directory, nil
}

// WriteToDisk writes the directory into disk
func (directory Directory) WriteToDisk(file *os.File) {

	// We write always at the beginning of file
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	var data = []interface{}{
		directory.Metadata,
		directory.FreeDataBlockArray,
		directory.FCBArray,
	}

	// Storing values into buffer
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed: ", err)
		}
	}

	// Dumping the buffer into disk
	file.Write(buf.Bytes())

}

func (directory Directory) String() string {
	return fmt.Sprintf("Directory Information\n%s", directory.Metadata)
}

// TestDirectory is a unit test for Directory
func TestDirectory() {
	file, err := os.OpenFile(".pfs", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	metadata := NewMetadata("pfs")
	directory := Directory{
		Metadata: metadata,
	}

	fmt.Println(directory)
	directory.WriteToDisk(file)
	directory, err = ReadDirectoryFromDisk(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(directory)
}
