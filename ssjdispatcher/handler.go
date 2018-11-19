package ssjdispatcher

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/service/sqs"
)

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
		bucket, err := GetValueFromDict(record.(map[string]interface{}), []string{"s3", "bucket", "name"})
		if err != nil {
			log.Println(err)
			continue
		}
		key, err := GetValueFromDict(record.(map[string]interface{}), []string{"s3", "object", "key"})
		if err != nil {
			log.Println(err)
			continue
		}
		bucketName := bucket.(string)
		keyName := key.(string)

		objectPath := "s3://" + bucketName + "/" + keyName

		for pattern, handleImage := range handler.mapHandler {
			re := regexp.MustCompile(pattern)
			if re.MatchString(objectPath) {
				result, err := createK8sJob(objectPath, handleImage)
				if err != nil {
					log.Println(err)
					return err
				}
				out, err := json.Marshal(result)
				if err != nil {
					log.Println(err)
					return err
				}
				log.Println(string(out))
			}

		}
	}
	return nil

}
