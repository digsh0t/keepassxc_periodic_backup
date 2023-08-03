package cliargs

import (
	"errors"
	"flag"
)

func GetArguments() (string, string, string, error) {
	pathPointer := flag.String("path", "", "Path to KeepassXC backup file")
	bucketPointer := flag.String("bucket", "", "S3 Bucket name")
	objectNamePointer := flag.String("object", "", "Object name on bucket")

	flag.Parse()
	if *pathPointer == "" {
		return "", "", "", errors.New("path should not be empty")
	}
	if *bucketPointer == "" {
		return "", "", "", errors.New("bucket name should not be empty")
	}
	if *objectNamePointer == "" {
		*objectNamePointer = *pathPointer
	}
	return *pathPointer, *bucketPointer, *objectNamePointer, nil
}
