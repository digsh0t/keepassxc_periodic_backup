package validate

import (
	"fmt"
	"os"
)

func ValidateKeepassXCBackupPath(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	fmt.Println(fi.Mode())
	return true, nil
}
