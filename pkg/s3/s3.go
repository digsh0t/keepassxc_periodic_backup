package s3

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/digsh0t/keepassxc_periodic_backup/pkg/utils"
)

type BucketBasics struct {
	S3Client *s3.Client
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func (basics BucketBasics) CreateBucket(name string, region string) error {
	_, err := basics.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		log.Printf("Couldn't create bucket %v in Region %v. Here's why: %v\n",
			name, region, err)
	}
	return err
}

func (s3client BucketBasics) GetBucket(bucketName string) (bool, error) {

	// Get the first page of results for ListObjectsV2 for a bucket
	_, err := s3client.S3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	}

	return exists, err
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (basics BucketBasics) UploadFile(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err := basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

func (bucket BucketBasics) GetFileEtag(bucketName string, objectKey string) (string, error) {
	objectResult, err := bucket.S3Client.GetObjectAttributes(context.TODO(), &s3.GetObjectAttributesInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		ObjectAttributes: []types.ObjectAttributes{
			types.ObjectAttributesChecksum,
			types.ObjectAttributesObjectParts,
			types.ObjectAttributesEtag,
		},
	})
	if err != nil {
		return "", err
	}
	return *(objectResult.ETag), nil
}

func (bucket BucketBasics) KeyExists(bucketName string, objectKey string) (bool, error) {
	_, err := bucket.S3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (bucket BucketBasics) CheckObjectUpToDate(bucketName string, objectKey string, path string) (bool, error) {
	// Get local file's MD5 hash
	hash, err := utils.GetFileMD5Hash(path)
	if err != nil {
		return false, err
	}
	hashString := fmt.Sprintf("%x", hash.Sum(nil))

	// Get S3 object Etag(MD5 hash)
	etag, err := bucket.GetFileEtag(bucketName, objectKey)
	if err != nil {
		return false, err
	}

	// If two hashes equal, return the program
	if hashString == etag {
		return true, nil
	}
	return false, nil
}
