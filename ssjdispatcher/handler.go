package ssjdispatcher

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/service/sqs"
)

const INDEXING string = "indexing"
const USERSYNC string = "usersync"

type S3ObjectHandler struct {
	mapHandler map[string]string
}

func (handler *S3ObjectHandler) InitHandlerMap() {
	handler.mapHandler = make(map[string]string)
}

func (handler *S3ObjectHandler) AddHandlerWithPattern(pattern string, jobName string) {
	handler.mapHandler[pattern] = jobName
}

// getBucketName returns bucket name of s3 object
func (handler *S3ObjectHandler) getBucketName(message *sqs.Message) (string, error) {
	return "", nil
}

// getKeyName returns s3 object key
func (handler *S3ObjectHandler) getKeyName(message *sqs.Message) (string, error) {
	return "", nil
}

// HandleS3Object handle s3 object
func (handler *S3ObjectHandler) HandleS3Object(message *sqs.Message) error {
	objectPath, err := handler.getS3ObjectPath(message)
	if err != nil {
		return fmt.Errorf("S3ObjectHandler: Can not get the object path of %s. Details %s", message, err)
	}

	for pattern, handleImage := range handler.mapHandler {
		re := regexp.MustCompile(pattern)
		if re.MatchString(objectPath) {
			_, err := createK8sJob(objectPath, handleImage)
			return err
		}

	}
	return nil

}

// getS3ObjectPath gets S3 object path
func (handler *S3ObjectHandler) getS3ObjectPath(message *sqs.Message) (string, error) {
	bucket, err := handler.getBucketName(message)
	if err != nil {
		return "", nil
	}
	key, err := handler.getKeyName(message)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	return "s3://" + bucket + "/" + key, nil
}
