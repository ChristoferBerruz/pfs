package pfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// FCB represents a logical, in RAM, File Control Block.
// We use byte arrays to easily read and write to/from the pfd
type FCB struct {
	FileName        [fileNameSize]byte
	CreateDateTime  [createDateTimeSize]byte
	FileSize        uint16
	StartingBlockID uint8
	Remarks         [remarksSize]byte
	CanBeWritten    bool
}

func invalidSizeMessage(variableName string, size int) string {
	return fmt.Sprintf("%s should be <= than %d in bytes", variableName, size)
}

// NewFCB returns a new FCB struct taking the datetime as now()
func NewFCB(fileName string, fileSize uint16, startingBlockID uint8) (FCB, error) {

	if fileNameSize-len(fileName) < 0 {
		return FCB{}, errors.New(invalidSizeMessage("FileName", fileNameSize))
	}

	block := FCB{
		FileSize:        fileSize,
		StartingBlockID: startingBlockID,
		CanBeWritten:    false,
	}

	// Populating the []byte for FileName
	for i := 0; i < len(fileName); i++ {
		block.FileName[i] = fileName[i]
	}

	// Populating the []byte for CreateDateTime
	dateString := time.Now().Format(dateTimeFormat)
	for i := 0; i < len(dateString); i++ {
		block.CreateDateTime[i] = dateString[i]
	}

	return block, nil
}

// ReadFCBFromDisk reads bytes from file at offset and returns an instance of FCB
func ReadFCBFromDisk(file *os.File, offset int64) (FCB, error) {
	if _, err := file.Seek(offset, 0); err != nil {
		return FCB{}, err
	}

	// Reading information to a buffer
	buf := make([]byte, fcbDiskSize)
	file.Read(buf)

	// block to be returned
	var block FCB
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, &block); err != nil {
		return FCB{}, err
	}
	return block, nil
}

// WriteToDisk writes a logcal FCB to the filesystem
func (block FCB) WriteToDisk(file *os.File, offset int64) {

	if _, err := file.Seek(offset, 0); err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	var data = []interface{}{
		block.FileName,
		block.CreateDateTime,
		block.FileSize,
		block.StartingBlockID,
		block.Remarks,
		block.CanBeWritten,
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

func (block FCB) String() string {
	return fmt.Sprintf("%s %2s %d %d %s", block.FileName, block.CreateDateTime,
		block.FileSize, block.StartingBlockID, block.Remarks)
}

// ModifyRemarks adds remarks in FCB in place
func (block *FCB) ModifyRemarks(remarks string) error {
	if remarksSize-len(remarks) < 0 {
		return errors.New(invalidSizeMessage("Remarks", remarksSize))
	}
	for i := 0; i < len(remarks); i++ {
		block.Remarks[i] = remarks[i]
	}
	return nil
}

// TestFCB is a simple test for FCB
func TestFCB() {
	file, err := os.OpenFile(".pfs", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	block, err := NewFCB("test.txt", 12, 3)
	if err != nil {
		fmt.Println(err)
	}
	block.WriteToDisk(file, 0)
	block1, _ := ReadFCBFromDisk(file, 0)
	fmt.Println(block1)
	block.WriteToDisk(file, 64)
	block2, _ := ReadFCBFromDisk(file, 64)
	fmt.Println(block2)
}
