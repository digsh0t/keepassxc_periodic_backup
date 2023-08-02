package main

import (
	"fmt"
	"log"

	"github.com/digsh0t/keepassxc_periodic_backup/pkg/cliargs"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/validate"
)

func main() {
	// TODO: Read input from users for back up file path, else use default file path
	path, err := cliargs.GetArguments()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return
	}
	fmt.Println(path)

	// TODO: Check if backup file is available and legit
	isLegit, err := validate.ValidateKeepassXCBackupPath(path)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	fmt.Println(isLegit)

	// TODO: Check if bucket is already existed on S3

	// TODO: A function that create new S3 Bucket from a bucket input
}
