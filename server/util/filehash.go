package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hashFunc := sha256.New()
	if _, err := io.Copy(hashFunc, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hashFunc.Sum(nil)), nil
}
