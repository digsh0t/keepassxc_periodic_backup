package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/cliargs"
	s3module "github.com/digsh0t/keepassxc_periodic_backup/pkg/s3"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/terraform"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/utils"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/validate"
)

func main() {
	// Read input from users for back up file path, else use default file path
	path, bucketName, objectName, err := cliargs.GetArguments()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return
	}

	// Check if backup file is available and legit
	isLegit, err := validate.ValidateKeepassXCBackupPath(path)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if !isLegit {
		log.Fatalf("ERROR: %s", "Not a KeepassXC Backup File")
	}

	// Check if bucket is already existed on S3
	exists, err := s3module.GetBucket("keepassxc-backup-bucket")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if !exists {
		log.Fatalf("ERROR: %s", "Bucket not existed")
	}

	terraform.ApplyS3Bucket(bucketName)

	// Initialize new SDK Client
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Fatalf("ERROR: %s", err)
	}

	// Get local file's MD5 hash
	hash, err := utils.GetFileMD5Hash(path)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	hashString := fmt.Sprintf("%x", hash.Sum(nil))

	// Get S3 Client
	bucket := s3module.BucketBasics{S3Client: s3.NewFromConfig(sdkConfig)}

	// Get S3 object Etag(MD5 hash)
	etag, err := bucket.GetFileEtag(bucketName, objectName)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	// If two hashes equal, return the program
	if hashString == etag {
		log.Println("Object up to date, no need to upload")
		return
	}

	// Upload the local file to S3
	err = bucket.UploadFile(bucketName, objectName, path)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("Uploaded %s successfully\n", path)
}
