package ssjdispatcher

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// S3Credentials contains AWS credentials
type S3Credentials struct {
	region             string
	awsAccessKeyID     string
	awsSecretAccessKey string
}

type AwsClient struct {
	credentials S3Credentials
	session     *session.Session
}

// LoadCredentialFromConfigFile loads AWS credentials from the config file
func (client *AwsClient) LoadCredentialFromConfigFile(path string) error {
	// Read data file
	jsonBytes, err := ReadFile(path)
	if err != nil {
		return err
	}

	// Get AWS region
	data, err := GetValueFromKeys(jsonBytes, []string{"AWS", "region"})
	if err != nil {
		return err
	}
	client.credentials.region = data.(string)

	// Get AWS access key id
	data, err = GetValueFromKeys(jsonBytes, []string{"AWS", "aws_access_key_id"})
	if err != nil {
		return err
	}
	client.credentials.awsAccessKeyID = data.(string)

	// Get AWS secret access key
	data, err = GetValueFromKeys(jsonBytes, []string{"AWS", "aws_secret_access_key"})
	if err != nil {
		return err
	}
	client.credentials.awsSecretAccessKey = data.(string)

	return nil
}

// createNewSession creats a aws s3 session
func (client *AwsClient) CreateNewSession() error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(client.credentials.region),
		Credentials: credentials.NewStaticCredentials(
			client.credentials.awsAccessKeyID, client.credentials.awsSecretAccessKey, ""),
	})
	client.session = sess

	return err
}

func (client *AwsClient) GetClientSession() *session.Session {
	return client.session
}
