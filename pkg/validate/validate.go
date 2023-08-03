package validate

import (
	"os"
	"path/filepath"
)

func ValidateKeepassXCBackupPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	ext := filepath.Ext(path)
	if ext != ".kdbx" {
		return false, nil
	}
	return true, nil
}
