package pfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
)

// Sizes are number of bytes
const (
	dateTimeFormat     string = "2006-01-02 15:04:05"
	createDateTimeSize        = 19
	fileNameSize              = 20
	remarksSize               = 22
	fcbDiskSize               = 64
)

// FCB represents a logical, in RAM, File Control Block.
// We use byte arrays to easily read and write to/from the pfd
type FCB struct {
	FileName        [fileNameSize]byte
	CreateDateTime  [createDateTimeSize]byte
	FileSize        uint16
	StartingBlockID uint8
	Remarks         [remarksSize]byte
}

// NewFCB returns a new FCB struct taking the datetime as now()
func NewFCB(fileName string, fileSize uint16, startingBlockID uint8) FCB {

	block := FCB{
		FileSize:        fileSize,
		StartingBlockID: startingBlockID,
	}
	// Transforming fileName to []byte
	var fnameBuf [len(block.FileName)]byte
	for i := 0; i < len(fileName); i++ {
		fnameBuf[i] = fileName[i]
	}

	// Creating a timestamp and storing it into []byte
	var dateTimeBuf [len(block.CreateDateTime)]byte
	dateString := time.Now().Format(dateTimeFormat)
	for i := 0; i < len(dateString); i++ {
		dateTimeBuf[i] = dateString[i]
	}

	var remarks [len(block.Remarks)]byte

	block.FileName = fnameBuf
	block.CreateDateTime = dateTimeBuf
	block.Remarks = remarks
	return block
}

// ReadFCBFromDisk reads bytes from file at offset and returns an instance of FCB
func ReadFCBFromDisk(file *os.File, offset int64) FCB {
	if _, err := file.Seek(offset, 0); err != nil {
		log.Fatal(err)
	}

	// Reading information to a buffer
	buf := make([]byte, fcbDiskSize)
	file.Read(buf)

	// block to be returned
	var block FCB
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, &block); err != nil {
		fmt.Println("Binary read failed: ", err)
	}
	return block
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

// TestFCB is a simple test for FCB
func TestFCB() {
	file, _ := os.OpenFile(".pfs", os.O_RDWR, 0644)
	block := NewFCB("test.txt", 12, 3)
	block.WriteToDisk(file, 0)
	block3 := NewFCB("del.txt", 17, 5)
	block3.WriteToDisk(file, 64)
	block2 := ReadFCBFromDisk(file, 0)
	block4 := ReadFCBFromDisk(file, 64)
	fmt.Println(block2)
	fmt.Println(block4)
}
