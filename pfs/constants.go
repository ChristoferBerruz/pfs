package pfs

// All sizes are number of bytes

// Overall design constants
const (
	volumeSize      = 10000
	totalDataBlocks = 31
	dataBlockSize   = 256
	dataAddress     = 2064
	directorySize   = 2038
)

// FCB constants
const (
	dateTimeFormat     string = "2006-01-02 15:04:05"
	createDateTimeSize        = 19
	fileNameSize              = 20
	remarksSize               = 21
	fcbDiskSize               = 64
)

// Constants for directory
const (
	nameOfVolumeSize = 20
)
