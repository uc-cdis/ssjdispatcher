package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/glog"
	mq "github.com/remind101/mq-go"
)

type SQSHandler struct {
	QueueURL   string
	Start      bool
	JobConfigs []JobConfig
	Server     *mq.Server
}

type JobConfig struct {
	Name        string      `name`
	Pattern     string      `pattern`
	Image       string      `image`
	ImageConfig interface{} `image_config`
}

// NewSQSHandler creates new SQSHandler instance
func NewSQSHandler(queueURL string) *SQSHandler {
	sqsHandler := new(SQSHandler)
	sqsHandler.QueueURL = queueURL
	//sqsHandler.PatternMap = GetNewImagePatternMap()
	return sqsHandler
}

// StartServer starts a server
func (handler *SQSHandler) StartServer() error {
	// return nil if the server already start
	if handler.Server != nil {
		return nil
	}

	fmt.Println("Start a new server")
	newClient, err := NewSQSClient()
	if err != nil {
		return err
	}

	handler.Server = mq.NewServer(handler.QueueURL, mq.HandlerFunc(func(m *mq.Message) error {
		return handler.HandleSQSMessage(m)
	}), mq.WithClient(newClient))
	handler.Server.Start()

	return nil

}

// ShutdownServer shutdowns a server
func (handler *SQSHandler) ShutdownServer() error {
	fmt.Println("Shutdown the server")
	if handler.Server == nil {
		return nil
	}
	err := handler.Server.Shutdown(context.Background())
	handler.Server = nil
	return err
}

/*
getObjectFromSQSMessage returns s3 object from sqs message

The format of a SQS message body:
{
	"Type" : "Notification",
  	"MessageId" : "f0207b9c-7255-5f61-998a-2f0e64c6eef0",
	"TopicArn" : "arn:aws:sns:us-east-1:707,
	"Subject" : "Amazon S3 Notification",
	"Message":  {"Records":[
		{"eventVersion":"2.0","eventSource":"aws:s3","awsRegion":"us-east-1",
		"eventTime":"2018-11-19T00:57:57.693Z","eventName":"ObjectCreated:Put",
		"userIdentity":{"principalId":"AWS:AIDAIU3LRUEK5OHS6FXRQ"},
		"requestParameters":{"sourceIPAddress":"173.24.34.163"},
		"responseElements":{"x-amz-request-id":"91CF670A054E0332",
		"x-amz-id-2":"h0ZQgg6w2qzKUkzivRizP1E1Jf9QAXSu1bUllWaF2b7j/63XRgjGLMNyI7sl016QKSaxK1L2RrI="},
		"s3":{"s3SchemaVersion":"1.0","configurationId":"Giang Bui",
		"bucket":{"name":"xssxs","ownerIdentity":{"principalId":"A365FU9MXCCF0S"},
		"arn":"arn:aws:s3:::xssxs"},"object":{"key":"api.py","size":8005,"eTag":"b4ef93035ff791f7d507a47342c89cd6",
		"sequencer":"005BF20A95A51A4C46"}}}]}
	}
}
*/

func getObjectsFromSQSMessage(m *mq.Message) []string {
	objectPaths := make([]string, 0)
	mapping := make(map[string][]interface{})
	msgBody := aws.StringValue(m.SQSMessage.Body)

	msgBodyInf, err := GetValueFromJSON([]byte(msgBody), []string{"Message"})
	if err != nil {
		glog.Infoln("The message is not the one from the bucket POST/PUT events. Detail ", err)
		return objectPaths
	}

	msgBody2 := msgBodyInf.(string)
	if err := json.Unmarshal([]byte(msgBody2), &mapping); err != nil {
		glog.Infoln("The message is not the one from the bucket POST/PUT events. Detail ", err)
		return objectPaths
	}

	records := mapping["Records"]
	for _, record := range records {
		recordByte, err := json.Marshal(record)
		if err != nil {
			glog.Errorln(err)
			continue
		}
		bucket, err := GetValueFromJSON(recordByte, []string{"s3", "bucket", "name"})
		if err != nil {
			glog.Errorln(err)
			continue
		}
		key, err := GetValueFromJSON(recordByte, []string{"s3", "object", "key"})
		if err != nil {
			glog.Errorln(err)
			continue
		}
		bucketName := bucket.(string)
		keyName := key.(string)

		objectPaths = append(objectPaths, "s3://"+bucketName+"/"+keyName)
	}

	return objectPaths
}

/*
HandleSQSMessage handles SQS message

The function takes a sqs message as input, extract the object urls and
decide which image need to be pulled to handle the s3 object
based on the object url

A SQS message may contains multiple records. The service goes though all
the records and compute the number of records need to be processed base
on their url and jobConfig list. If the number is larger than the availbility
of jobpool, the service will take a sleep until the resource is available.

If the function returns an error other than nil, the message is put back
to the queue and retry later (handled by `md` library). That makes sure
the message is properly handle before it actually deleted

*/
func (handler *SQSHandler) HandleSQSMessage(m *mq.Message) error {
	objectPaths := getObjectsFromSQSMessage(m)

	jobNameList := make([]string, 0)
	for _, jobConfig := range handler.JobConfigs {
		jobNameList = append(jobNameList, jobConfig.Name)
	}

	// remove completed jobs
	RemoveCompletedJobs(jobNameList)

	jobMap := make(map[string]JobConfig)
	for _, objectPath := range objectPaths {
		for _, jobConfig := range handler.JobConfigs {
			re := regexp.MustCompile(jobConfig.Pattern)
			if re.MatchString(objectPath) {
				jobMap[objectPath] = jobConfig
			}
		}
	}

	for len(jobMap)+GetNumberRunningJobs(jobNameList) > GetMaxJobConfig() {
		time.Sleep(5 * time.Second)
	}

	glog.Infof("Start to run %d jobs", len(jobMap))

	for objectPath, jobConfig := range jobMap {
		glog.Info("Processing: ", objectPath)
		result, err := CreateK8sJob(objectPath, jobConfig)
		if err != nil {
			glog.Errorln(err)
			return err
		}
		out, err := json.Marshal(result)
		if err != nil {
			glog.Errorln(err)
			return err
		}
		glog.Info(string(out))
	}

	return nil
}

func (handler *SQSHandler) handleAddNewJobConfig(jsonBytes []byte) error {
	jobConfig := JobConfig{}
	if err := json.Unmarshal(jsonBytes, &jobConfig); err != nil {
		return err
	}
	if jobConfig.Name != "" && jobConfig.Image != "" {
		handler.JobConfigs = append(handler.JobConfigs, jobConfig)
	} else {

		return errors.New("Name and image args are required")
	}
	return nil
}

func (handler *SQSHandler) handleDeleteJobConfig(pattern string) error {
	for idx, job := range handler.JobConfigs {
		if job.Pattern == pattern {
			handler.JobConfigs = append(handler.JobConfigs[:idx], handler.JobConfigs[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("There is no job with provided pattern\n %s", pattern)
}

func (handler *SQSHandler) handleListJobConfigs() (string, error) {
	str := ""
	for _, job := range handler.JobConfigs {
		jsonBytes, err := json.Marshal(job)
		if err != nil {
			return "", err
		}
		str = str + string(jsonBytes) + ","
	}
	return "[" + str + "]", nil
}
