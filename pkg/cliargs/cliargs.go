package cliargs

import (
	"errors"
	"flag"
)

func GetArguments() (string, error) {
	pathPointer := flag.String("path", "", "Path to KeepassXC backup file")

	flag.Parse()
	if *pathPointer == "" {
		return "", errors.New("path should not be empty")
	}
	return *pathPointer, nil
}
