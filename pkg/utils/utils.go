package utils

import (
	"crypto/md5"
	"hash"
	"io"
	"os"
)

func GetFileMD5Hash(filepath string) (hash.Hash, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return nil, err
	}
	return hash, nil
}
