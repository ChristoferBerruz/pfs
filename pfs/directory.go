package pfs

const (
	directorySize    = 1040          // 4*blocks of 256 bytes + 16 bytes
	dataStartAddress = directorySize //Data block address starts at 1040
	totalDataBlocks  = 35
)

// Metadata stores properties about the volume
type Metadata struct {
	directorySize       uint16
	numberOfFilesInDisk uint8
	availableFCBPointer uint16 // offset to start writting FCBs
}

// Directory is a directory of FileSystem
type Directory struct {
	metadata   Metadata
	freeBlocks [totalDataBlocks]byte
	fcbList    []FCB
}
