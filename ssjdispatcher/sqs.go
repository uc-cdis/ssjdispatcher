package ssjdispatcher

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// the path to aws credential
const CREDPATH = "./credentials.json"

// getSQSSession gets new sqs session
func GetSQSSession() (*sqs.SQS, error) {
	awsClient := AwsClient{}
	awsClient.LoadCredentialFromConfigFile(CREDPATH)
	awsClient.CreateNewSession()

	return sqs.New(awsClient.GetClientSession()), nil
}

// getSQSMessages gets messages from queue
func GetSQSMessages(svc *sqs.SQS, queueUrl string) (*sqs.ReceiveMessageOutput, error) {
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

func DeleteSQSMessage(svc *sqs.SQS, queueUrl string, sqsMessage *sqs.Message) error {
	// // Delete message
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),     // Required
		ReceiptHandle: sqsMessage.ReceiptHandle, // Required
	}
	_, err := svc.DeleteMessage(deleteParams) // No response returned when successed.
	return err
}
