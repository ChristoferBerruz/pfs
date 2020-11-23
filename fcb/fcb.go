package fcb

import (
	"fmt"
	"time"
)

// FCB represents a File Control Block
type FCB struct {
	fileName        string
	fileSize        uint32
	createDateTime  time.Time
	startingBlockID uint32
	endingBlockID   uint32
}

// NewFCB creates a new FCB struct using local datetime
func NewFCB(fileName string, fileSize uint32, startingBlockID uint32, endingBlockID uint32) FCB {
	return FCB{
		fileName:        fileName,
		fileSize:        fileSize,
		startingBlockID: startingBlockID,
		endingBlockID:   endingBlockID,
		createDateTime:  time.Now(),
	}
}

func (fcb FCB) String() string {
	return fmt.Sprintf("%s %5d %5d %d", fcb.fileName,
		fcb.fileSize, fcb.startingBlockID, fcb.endingBlockID)
}
