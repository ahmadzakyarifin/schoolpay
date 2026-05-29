package utils

import (
	"os"
)

func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0777)
	}
	return nil
}
