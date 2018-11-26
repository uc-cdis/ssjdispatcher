package ssjdispatcher

import (
	"encoding/json"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSHandler struct {
	QueueURL   string
	Start      bool
	PatternMap *ImagePatternMap
}

// NewSQSHandler creates new SQSHandler instance
func NewSQSHandler(queueURL string, start bool) *SQSHandler {
	sqsHandler := new(SQSHandler)
	sqsHandler.QueueURL = queueURL
	sqsHandler.Start = start
	sqsHandler.PatternMap = GetNewImagePatternMap()
	return sqsHandler
}

// HandleSQSMessage handles SQS messages
func (handler *SQSHandler) HandleSQSMessage(message *sqs.Message) error {
	mapping := make(map[string][]interface{})
	msgBody := *message.Body
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

// StartSQSQuerying periodically queries SQS URL to get messages
func StartSQSQuerying(sqsHandler *SQSHandler) {
	svc, err := GetSQSSession()
	if err != nil {
		panic(err)
	}
	for {
		if sqsHandler.Start == false {
			time.Sleep(10 * time.Second)
			log.Println("The service is sleeping")
			continue
		}

		time.Sleep(2 * time.Second)
		log.Println("The service is running")

		resp, err := GetSQSMessages(svc, sqsHandler.QueueURL)
		if err != nil {
			log.Printf("Can not query %s. Please check the URL. Detail %s\n", sqsHandler.QueueURL, err)
		}

		for _, message := range resp.Messages {
			//log.Printf("[Receive message] \n%v \n\n", message)
			err := sqsHandler.HandleSQSMessage(message)
			if err != nil {
				log.Printf("Can not handle message \n%v \n\nDetail: %s", message, err)
				continue
			} else {
				if DeleteSQSMessage(svc, sqsHandler.QueueURL, message) != nil {
				} else {
					log.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
				}
			}
		}
	}
}
