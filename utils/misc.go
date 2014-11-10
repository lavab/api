package utils

import (
	"os"
)

// FileExists is a stupid little wrapper of os.Stat that checks whether a file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}
