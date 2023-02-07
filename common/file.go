package common

import (
	"os"
)

func GetFileContent(fileFullPath string) ([]byte, error) {
	return os.ReadFile(fileFullPath)
}
