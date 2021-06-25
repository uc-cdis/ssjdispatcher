package handlers

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	batchtypev1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

type mockSQSClient struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (m *mockSQSClient) SendMessage(opts *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := m.Called(opts)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

type mockJobClient struct {
	batchtypev1.JobInterface
	mock.Mock
}

func (m *mockJobClient) List(opts metav1.ListOptions) (*batchv1.JobList, error) {
	args := m.Called(opts)
	return args.Get(0).(*batchv1.JobList), args.Error(1)
}

func TestCheckJobs(t *testing.T) {
	sqsClient := &mockSQSClient{}
	sqsClient.On("SendMessage", mock.Anything).Return(&sqs.SendMessageOutput{
		MessageId: aws.String("new message id"),
	}, nil)
	jobClient := &mockJobClient{}
	jobClient.On("List", mock.Anything).Return(&batchv1.JobList{
		Items: []batchv1.Job{
			{
				Status: batchv1.JobStatus{
					Failed: 1,
				},
				ObjectMeta: metav1.ObjectMeta{
					UID: types.UID("job1"),
				},
			},
		},
	}, nil)

	messageBytes, _ := json.Marshal(map[string]string{
		"Message": string("test message"),
	})
	message := string(messageBytes)

	handler := &SQSHandler{
		MonitoredJobs: []*JobInfo{
			{
				UID:     "job1",
				Retries: MAX_RETRIES - 1,
				SQSMessage: &sqs.Message{
					MessageId: aws.String("existing message id"),
					Body:      aws.String(message),
				},
			},
		},
		sqsClient: sqsClient,
		jobHandler: &jobHandler{
			jobClient: jobClient,
		},
	}
	handler.checkJobs()

	// This is not a good thing. We are going to send a new message with a new job
	// id and so we end up duplicating one message to many over
	sqsClient.AssertCalled(t, "SendMessage", &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(""),
	})
	assert.Equal(t, handler.MonitoredJobs, []*JobInfo{
		{
			UID:     "job1",
			Retries: MAX_RETRIES,
			SQSMessage: &sqs.Message{
				MessageId: aws.String("existing message id"),
				Body:      aws.String(message),
			},
		},
	})
}
