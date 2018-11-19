package ssjdispatcher

import (
	"encoding/json"
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

// HandleS3Object handle s3 object
func (handler *S3ObjectHandler) HandleS3Objects(message *sqs.Message) error {
	mapping := make(map[string][]interface{})
	msgBody := *message.Body
	if err := json.Unmarshal([]byte(msgBody), &mapping); err != nil {
		panic(err)
	}
	records := mapping["Records"]
	for _, record := range records {
		//s3aw := record["s3"]
		bucket := record.(map[string]interface{})["s3"].(map[string]interface{})["bucket"].(map[string]interface{})["name"].(string)
		key := record.(map[string]interface{})["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)

		objectPath := "s3://" + bucket + "/" + key

		for pattern, handleImage := range handler.mapHandler {
			re := regexp.MustCompile(pattern)
			if re.MatchString(objectPath) {
				_, err := createK8sJob(objectPath, handleImage)
				return err
			}

		}
	}
	return nil

}
