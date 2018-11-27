package ssjdispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/remind101/mq-go"
)

type SQSHandler struct {
	QueueURL   string
	Start      bool
	PatternMap *ImagePatternMap
	Server     *mq.Server
}

// NewSQSHandler creates new SQSHandler instance
func NewSQSHandler(queueURL string, start bool) *SQSHandler {
	sqsHandler := new(SQSHandler)
	sqsHandler.QueueURL = queueURL
	sqsHandler.PatternMap = GetNewImagePatternMap()
	return sqsHandler
}

func (handler *SQSHandler) StartServer() {
	fmt.Println("Start a new server")
	if handler.Server == nil {
		handler.Server = mq.NewServer(handler.QueueURL, mq.HandlerFunc(func(m *mq.Message) error {
			return handler.HandleSQSMessage(m)
		}))
		handler.Server.Start()
	}

}

func (handler *SQSHandler) ShutdownServer() {
	fmt.Println("Shutdown the server")
	if handler.Server == nil {
		return
	}
	handler.Server.Shutdown(context.Background())
	handler.Server = nil
}

func (handler *SQSHandler) HandleSQSMessage(m *mq.Message) error {
	mapping := make(map[string][]interface{})
	msgBody := aws.StringValue(m.SQSMessage.Body)
	if err := json.Unmarshal([]byte(msgBody), &mapping); err != nil {
		return err
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

		fmt.Println("Processing: ", objectPath)

		for pattern, handleImage := range handler.PatternMap.Mapping {
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
