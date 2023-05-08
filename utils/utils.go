package utils

import (
	"os"
)

func PathExists(modelPath string) bool {
	_, err := os.Stat(modelPath)
	return err == nil
}
