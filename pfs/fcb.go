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
	dateTimeFormat string = "2006-01-02 15:04:05"
	fcbDiskSize           = 63
)

// FCB represents a logical, in RAM, File Control Block.
// We use byte arrays to easily read and write to/from the pfd
type FCB struct {
	FileName        [20]byte
	CreateDateTime  [19]byte
	FileSize        uint64
	StartingBlockID uint64
	EndingBlockID   uint64
}

// NewFCB returns a new FCB struct taking the datetime as now()
func NewFCB(fileName string, fileSize, startingBlockID, endingBlockID uint64) FCB {

	// Transforming fileName to []byte
	var fnameBuf [20]byte
	for i := 0; i < len(fileName); i++ {
		fnameBuf[i] = fileName[i]
	}

	// Creating a timestamp and storing it into []byte
	var dateTimeBuf [19]byte
	dateString := time.Now().Format(dateTimeFormat)
	for i := 0; i < len(dateString); i++ {
		dateTimeBuf[i] = dateString[i]
	}

	return FCB{
		FileName:        fnameBuf,
		CreateDateTime:  dateTimeBuf,
		FileSize:        fileSize,
		StartingBlockID: startingBlockID,
		EndingBlockID:   endingBlockID,
	}
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
		block.EndingBlockID,
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
