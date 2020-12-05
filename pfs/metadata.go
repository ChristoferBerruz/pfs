package pfs

import (
	"errors"
	"fmt"
)

// Metadata contains a summary of the volume
type Metadata struct {
	NameOfVolume        [nameOfVolumeSize]byte
	DataInitialAddress  uint16
	NumberOfFilesStored uint8
}

// NewMetadata returns a new metadata instance
func NewMetadata(nameOfVolume string) Metadata {

	metadata := Metadata{}
	err := (&metadata).setNameOfVolume(nameOfVolume)
	if err != nil {
		fmt.Println(err)
	}
	metadata.DataInitialAddress = dataAddress
	metadata.NumberOfFilesStored = 0
	return metadata
}

func (metadata *Metadata) setNameOfVolume(volName string) error {
	if nameOfVolumeSize-len(volName) < 0 {
		return errors.New(invalidSizeMessage("VolumeName", nameOfVolumeSize))
	}

	for i := 0; i < len(volName); i++ {
		(*metadata).NameOfVolume[i] = volName[i]
	}
	return nil
}

func (metadata Metadata) getNameOfVolume() string {
	return fmt.Sprintf("%s", metadata.NameOfVolume)
}
func (metadata Metadata) String() string {
	return fmt.Sprintf("Name of volume: %s\nDataInitialAddress:%d\nFiles Stored:%d\n",
		metadata.NameOfVolume, metadata.DataInitialAddress, metadata.NumberOfFilesStored)
}
