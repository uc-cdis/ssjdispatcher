package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/uc-cdis/ssjdispatcher/ssjdispatcher"
)

// the path to aws credential
const CREDPATH = "./credentials.json"

// const (
// 	QueueUrl = "https://sqs.us-east-1.amazonaws.com/440721843528/mySQS"
// )

// getSQSSession gets new sqs session
func getSQSSession() (*sqs.SQS, error) {
	awsClient := ssjdispatcher.AwsClient{}
	awsClient.LoadCredentialFromConfigFile(CREDPATH)
	awsClient.CreateNewSession()

	return sqs.New(awsClient.GetClientSession()), nil
}

// getSQSMessages gets messages from queue
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

// go run main.go https://sqs.us-east-1.amazonaws.com/440721843528/mySQS '{"pattern1": "do_1", "pattern2": "do_2"}'
func main() {

	argsWithProg := os.Args
	if len(argsWithProg) < 3 {
		panic("The service requires 3. Only " + strconv.Itoa((len(argsWithProg))) + " provided")
	}

	queueURL := argsWithProg[1]
	mappingStr := argsWithProg[2]

	//mappingStr = "{\"pattern1\": \"do_1\", \"pattern2\": \"do_2\"}"
	//queueURL = "https://sqs.us-east-1.amazonaws.com/440721843528/mySQS"

	var mapping map[string]string
	if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
		panic(err)
	}

	objectHandler := ssjdispatcher.S3ObjectHandler{}
	objectHandler.InitHandlerMap()
	for pattern, quayImage := range mapping {
		objectHandler.AddHandlerWithPattern(pattern, quayImage)
	}

	svc, err := getSQSSession()
	if err != nil {
		panic(err)
	}
	resp, _ := getSQSMessages(svc, queueURL)

	for _, message := range resp.Messages {
		// fmt.Printf("[Receive message] \n%v \n\n", message)
		err := objectHandler.HandleS3Object(message)
		if err != nil {
			log.Println(err)
		} else {
			if deleteSQSMessage(svc, queueURL, message) != nil {
			} else {
				log.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
			}
		}
	}

}
