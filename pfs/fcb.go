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
	FileName          [fileNameSize]byte
	CreateDateTime    [createDateTimeSize]byte
	FileSize          uint16
	StartingBlockID   uint8
	Remarks           [remarksSize]byte
	ContainsValidData bool
}

func (fcb *FCB) setFileName(fileName string) error {
	if fileNameSize-len(fileName) < 0 {
		return errors.New(invalidSizeMessage("FileName", fileNameSize))
	}

	for i := 0; i < len(fileName); i++ {
		(*fcb).FileName[i] = fileName[i]
	}

	return nil
}

func (fcb FCB) getFileName() string {
	b := make([]byte, fileNameSize)

	for i := 0; i < fileNameSize; i++ {
		b[i] = fcb.FileName[i]
	}
	return string(b)
}

func (fcb FCB) getCreateDateTime() string {
	b := make([]byte, createDateTimeSize)

	for i := 0; i < createDateTimeSize; i++ {
		b[i] = fcb.CreateDateTime[i]
	}

	return string(b)
}

func (fcb *FCB) setCreateDateTime(dateTime string) error {
	if createDateTimeSize-len(dateTime) < 0 {
		return errors.New(invalidSizeMessage("CreateDateTime", createDateTimeSize))
	}

	for i := 0; i < len(dateTime); i++ {
		(*fcb).CreateDateTime[i] = dateTime[i]
	}

	return nil
}

func (fcb FCB) getRemarks() string {
	b := make([]byte, remarksSize)
	for i := 0; i < remarksSize; i++ {
		b[i] = fcb.Remarks[i]
	}

	return string(b)
}

func (fcb *FCB) setRemarks(remarks string) error {
	if remarksSize-len(remarks) < 0 {
		return errors.New(invalidSizeMessage("Remarks", remarksSize))
	}
	for i := 0; i < len(remarks); i++ {
		(*fcb).Remarks[i] = remarks[i]
	}
	return nil
}

// NewFCB returns a new FCB struct taking the datetime as now()
func NewFCB(fileName string, fileSize uint16, startingBlockID uint8) (FCB, error) {

	block := FCB{
		FileSize:          fileSize,
		StartingBlockID:   startingBlockID,
		ContainsValidData: true,
	}

	err := (&block).setFileName(fileName)
	if err != nil {
		return FCB{}, err
	}

	dateString := time.Now().Format(dateTimeFormat)
	(&block).setCreateDateTime(dateString)

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
func (fcb FCB) WriteToDisk(file *os.File, offset int64) {

	if _, err := file.Seek(offset, 0); err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	var data = []interface{}{
		fcb.FileName,
		fcb.CreateDateTime,
		fcb.FileSize,
		fcb.StartingBlockID,
		fcb.Remarks,
		fcb.ContainsValidData,
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

func (fcb FCB) String() string {
	return fmt.Sprintf("%s %2s %d %d %s\n", fcb.FileName, fcb.CreateDateTime,
		fcb.FileSize, fcb.StartingBlockID, fcb.Remarks)
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
