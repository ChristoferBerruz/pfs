package pfs

import "fmt"

func invalidSizeMessage(variableName string, size int) string {
	return fmt.Sprintf("%s should be <= than %d in bytes", variableName, size)
}
