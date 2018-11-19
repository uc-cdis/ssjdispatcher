package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/uc-cdis/ssjdispatcher/ssjdispatcher"
)

const CREDPATH = "./credentials.json"

// const (
// 	QueueUrl = "https://sqs.us-east-1.amazonaws.com/440721843528/mySQS"
// )

func getSQSSession() (*sqs.SQS, error) {
	awsClient := ssjdispatcher.AwsClient{}
	awsClient.LoadCredentialFromConfigFile(CREDPATH)
	awsClient.CreateNewSession()

	return sqs.New(awsClient.GetClientSession()), nil
}

func getSQSMessages(svc *sqs.SQS, queueUrl string) (*sqs.ReceiveMessageOutput, error) {
	// Receive message
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(20),
	}
	resp, err := svc.ReceiveMessage(receiveParams)
	return resp, err
}

func deleteSQSMessage(svc *sqs.SQS, queueUrl string, sqsMessage *sqs.Message) error {
	// // Delete message
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),     // Required
		ReceiptHandle: sqsMessage.ReceiptHandle, // Required
	}
	_, err := svc.DeleteMessage(deleteParams) // No response returned when successed.
	return err
	// if err != nil {
	// 	log.Println(err)
	// }
	//fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *sqsMessage.MessageId)

}

func main() {
	argsWithProg := os.Args
	queueURL := argsWithProg[1]
	var mapping map[string]string
	if err := json.Unmarshal([]byte(argsWithProg[2]), &mapping); err != nil {
		panic(err)
	}

	objectHandler := ssjdispatcher.S3ObjectHandler{}
	objectHandler.InitHandlerMap()
	for pattern, quayImage := range mapping {
		objectHandler.AddHandlerWithPattern(pattern, quayImage)
	}

	svc, err := getSQSSession()
	if err != nil {
		log.Printf("Can not get SQS session. Detail %s", err)
		panic("Can not get SQS session")
	}
	resp, _ := getSQSMessages(svc, queueURL)

	for _, message := range resp.Messages {
		// fmt.Printf("[Receive message] \n%v \n\n", message)
		err := objectHandler.HandleS3Object(message)
		if err != nil {
			log.Println(err)
		}
		deleteSQSMessage(svc, queueURL, message)

	}

}
