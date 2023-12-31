package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/cliargs"
	s3module "github.com/digsh0t/keepassxc_periodic_backup/pkg/s3"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/validate"
)

func main() {
	// Read input from users for back up file path, else use default file path
	path, bucketName, objectName, bucketRegion, err := cliargs.GetArguments()
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return
	}

	// Check if backup file is available and legit
	isLegit, err := validate.ValidateKeepassXCBackupPath(path)
	if err != nil {
		log.Fatalf("Failed to validate KeepassXC Backup Path with ERROR: %s", err)
	}
	if !isLegit {
		log.Fatalf("ERROR: %s", "Not a KeepassXC Backup File")
	}

	// Initialize new SDK Client
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(bucketRegion))
	if err != nil {
		log.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Fatalf("Failed to load the AWS SDK ERROR: %s", err)
	}

	// Get S3 Client
	s3client := s3module.BucketBasics{S3Client: s3.NewFromConfig(sdkConfig)}

	// Check if bucket is already existed on S3
	exists, err := s3client.GetBucket(bucketName)
	if err != nil {
		log.Fatalf("Failed to check bucket existent with ERROR: %s", err)
	}
	if !exists {
		err = s3client.CreateBucket(bucketName, bucketRegion)
		if err != nil {
			log.Fatalf("Failed to create S3 bucket with ERROR: %s", err)
		}
	}

	// Check if object existed on Bucket
	isExisted, err := s3client.KeyExists(bucketName, objectName)
	if err != nil {
		log.Fatalf("Failed to check object on S3 bucket with ERROR: %s", err)
	}
	if isExisted {
		isSame, err := s3client.CheckObjectUpToDate(bucketName, objectName, path)
		if err != nil {
			log.Fatalf("Failed to compare version on S3 bucket with ERROR: %s", err)
		}
		if isSame {
			log.Println("Object up to date, no need to upload")
			return
		}
	}

	// Upload the local file to S3
	err = s3client.UploadFile(bucketName, objectName, path)
	if err != nil {
		log.Fatalf("Failed to upload file to S3 with ERROR: %s", err)
	}
	log.Printf("Uploaded %s successfully\n", path)
}
