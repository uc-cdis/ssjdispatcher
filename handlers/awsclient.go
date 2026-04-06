package handlers

import (
	"context" // Required for all v2 API calls

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/golang/glog"
)

// S3Credentials contains AWS credentials
type AWSCredentials struct {
	region             string
	awsAccessKeyID     string
	awsSecretAccessKey string
}

// Custom SQSAPI interface to replace sqsiface.SQSAPI
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/migrate-gosdk.html#mocking-and-iface
type SQSAPI interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// NewSQSClient create new SQSAPI client
func NewSQSClient() (SQSAPI, error) {
	cred, err := loadCredentialFromConfigFile(LookupCredFile())
	if err != nil {
		return nil, err
	}

	var optFns []func(*config.LoadOptions) error

	if cred.region != "" {
		optFns = append(optFns, config.WithRegion(cred.region))
	}

	if cred.awsAccessKeyID != "" && cred.awsSecretAccessKey != "" {
		staticProvider := credentials.NewStaticCredentialsProvider(cred.awsAccessKeyID, cred.awsSecretAccessKey, "")
		optFns = append(optFns, config.WithCredentialsProvider(staticProvider))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), optFns...)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}

// loadCredentialFromConfigFile loads AWS credentials from the config file
func loadCredentialFromConfigFile(path string) (*AWSCredentials, error) {
	credentials := new(AWSCredentials)
	// Read data file
	jsonBytes, err := ReadFile(path)
	if err != nil {
		return nil, err
	}

	if region, err := GetValueFromJSON(jsonBytes, []string{"AWS", "region"}); err != nil {
		glog.Info("Can not read region from credential file. Gonna use attached service account")
	} else {
		credentials.region = region.(string)
	}

	if awsKeyID, err := GetValueFromJSON(jsonBytes, []string{"AWS", "aws_access_key_id"}); err != nil {
		glog.Info("Can not read aws key from credential file. Gonna use attached service account ")
	} else {
		credentials.awsAccessKeyID = awsKeyID.(string)
	}

	if awsSecret, err := GetValueFromJSON(jsonBytes, []string{"AWS", "aws_secret_access_key"}); err != nil {
		glog.Info("Can not read aws secret key from credential file. Gonna use attached service account")
	} else {
		credentials.awsSecretAccessKey = awsSecret.(string)
	}

	return credentials, nil
}
