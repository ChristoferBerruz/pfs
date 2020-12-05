package pfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

// FileSystem is the main struct to talk to pfs file system
type FileSystem struct {
	PfsFile    *os.File
	VolumeName string
	Directory  Directory
}

// OpenVolume opens the volume as a filesystem
func (fileSystem *FileSystem) OpenVolume(fileName string) {
	if !FileExists(fileName) {
		// We need to create file and initialize and empty pfs system inside of it
		err := CreateFile(fileName)
		if err != nil {
			fmt.Println(err)
		}

		directory := NewDirectory(fileName)

		file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}

		(*fileSystem).PfsFile = file
		(*fileSystem).VolumeName = fileName
		(*fileSystem).Directory = directory
	} else {

		file, _ := os.OpenFile(fileName, os.O_RDWR, 0644)
		directory, err := ReadDirectoryFromDisk(file)
		if err != nil {
			log.Fatal(err)
		}
		(*fileSystem).VolumeName = directory.Metadata.getNameOfVolume()
		(*fileSystem).PfsFile = file
		(*fileSystem).Directory = directory
	}

	fmt.Printf("Sucesfully opened %s volume\n", fileName)
	fmt.Println((*fileSystem).Directory)
}

// Quit quits the file system and flushes directory back to file
func (fileSystem *FileSystem) Quit() {

	err := fileSystem.flushDirectory()
	if err != nil {
		fmt.Println(err)
		return
	}

	(*fileSystem).PfsFile.Close()
	fmt.Println("Successfully exited file system")
}

// Kill deletes a given volume name if found
func (fileSystem *FileSystem) Kill(volName string) {
	if FileExists(volName) {
		if (*fileSystem).PfsFile != nil && (*fileSystem).VolumeName == volName {
			(*fileSystem).PfsFile.Close()
			(*fileSystem).PfsFile = nil
		}

		err := os.Remove(volName)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("Successfully kiled %s volume\n", volName)
}

// Put moves a file from host file system into pfs
func (fileSystem *FileSystem) Put(fileName string) {
	if (*fileSystem).PfsFile == nil {
		fmt.Println("You must open pfs first")
		return
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(data) > totalDataBlocks*dataBlockSize {
		fmt.Println("File is too big to be inserted in pfs")
		return
	}

	fileSize := uint16(len(data)) // Get the file size

	// Finding how many blocks is easy enough
	necesaryBlocks := uint16(math.Ceil(float64(fileSize) / float64(dataBlockSize)))
	foundBlocks := uint16(0)
	location := -1
	for idx, availability := range (*fileSystem).Directory.FreeDataBlockArray {
		if foundBlocks > 0 && !availability {
			foundBlocks = 0
			continue
		}

		if availability {
			foundBlocks++
		}
		if foundBlocks == necesaryBlocks {
			location = idx
			break
		}
	}

	if location == -1 {
		fmt.Println("Not enough free blocks to store file.")
		return
	}

	// We have found the location of data that should be written
	blockID := location - int(necesaryBlocks) + 1
	offset := dataAddress + blockID*dataBlockSize
	(*fileSystem).PfsFile.Seek(int64(offset), 0)
	(*fileSystem).PfsFile.Write(data)

	// Now we need to find a FCB to store the records of the value
	for idx, fcb := range (*fileSystem).Directory.FCBArray {
		if !fcb.ContainsValidData {
			block, _ := NewFCB(fileName, fileSize, uint8(blockID))
			(*fileSystem).Directory.FCBArray[idx] = block
			break
		}
	}

	// Now we also need to set the free data blocks as ocupied
	for i := 0; i < int(necesaryBlocks); i++ {
		(*fileSystem).Directory.FreeDataBlockArray[blockID+i] = false
	}

	// We know update the directory
	(*fileSystem).Directory.Metadata.NumberOfFilesStored++
}

// RemoveFile removes a file from file system
func (fileSystem *FileSystem) RemoveFile(fileName string) {

	lenOfName := len(fileName)

	for idx, fcb := range (*fileSystem).Directory.FCBArray {

		nameTruncated := string(fcb.getFileName()[0:lenOfName])

		if fcb.ContainsValidData && (nameTruncated == fileName) {
			// We found the file so we have to mark the fcb as invalid
			// AND the blocks of data as free
			// And modify metadata information
			(*fileSystem).Directory.FCBArray[idx].ContainsValidData = false

			// Setting the data blocks as free
			numberOfBlocks := uint8(math.Ceil(float64(fcb.FileSize) / float64(dataBlockSize)))
			for i := uint8(0); i < numberOfBlocks; i++ {
				(*fileSystem).Directory.FreeDataBlockArray[fcb.StartingBlockID+i] = true
			}

			// Set directory number of files in system
			(*fileSystem).Directory.Metadata.NumberOfFilesStored--
			fmt.Printf("Sucessfully removed %s from file system\n", fileName)
			return
		}
	}
	fmt.Println("File does not exist in filesystem. Nothing to delete.")

}

// PutRemarks modifies a FCB to set different remarks
func (fileSystem *FileSystem) PutRemarks(fileName string, remarks string) {

	lenOfName := len(fileName)

	for idx, fcb := range (*fileSystem).Directory.FCBArray {

		nameTruncated := string(fcb.getFileName()[0:lenOfName])

		if fcb.ContainsValidData && (nameTruncated == fileName) {
			err := (*fileSystem).Directory.FCBArray[idx].setRemarks(remarks)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Sucessfully modified remarks")
			return
		}
	}

	fmt.Println("Files does not exist in file something. Nothing to change")
}

// MoveOut moves fileName from .pfs to current directory
func (fileSystem *FileSystem) MoveOut(fileName string) {

	lenOfName := len(fileName)

	for _, fcb := range (*fileSystem).Directory.FCBArray {

		nameTruncated := string(fcb.getFileName()[0:lenOfName])

		if fcb.ContainsValidData && (nameTruncated == fileName) {
			blockID := fcb.StartingBlockID
			lenOfData := fcb.FileSize
			offset := dataAddress + uint16(blockID)*dataBlockSize

			buf := make([]byte, lenOfData)

			(*fileSystem).PfsFile.Seek(int64(offset), 0)

			(*fileSystem).PfsFile.Read(buf)

			// Create file
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Everything is good, so flush
			file.Write(buf)
			file.Close()

			// We moved it, so lets get' remove the file internally
			fileSystem.RemoveFile(fileName)
			return
		}
	}

	fmt.Println("Files does not exist in file something. Nothing to change")
}

// Dir shows files available in the file system
func (fileSystem *FileSystem) Dir() {
	result := ""
	for _, fcb := range (*fileSystem).Directory.FCBArray {
		if fcb.ContainsValidData {
			result += fcb.String()
		}
	}
	fmt.Println(result)
}

func (fileSystem *FileSystem) flushDirectory() error {
	if (*fileSystem).PfsFile == nil {
		return errors.New("Cannot flush to a nil file. Please open volume first")
	}

	(*fileSystem).Directory.WriteToDisk((*fileSystem).PfsFile)
	return nil
}
