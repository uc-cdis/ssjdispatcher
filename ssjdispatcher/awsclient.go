package ssjdispatcher

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// S3Credentials contains AWS credentials
type AWSCredentials struct {
	region             string
	awsAccessKeyID     string
	awsSecretAccessKey string
}

// loadCredentialFromConfigFile loads AWS credentials from the config file
func loadCredentialFromConfigFile(path string) (*AWSCredentials, error) {
	credentials := new(AWSCredentials)
	// Read data file
	jsonBytes, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	var mapping map[string]interface{}
	json.Unmarshal(jsonBytes, &mapping)
	if region, err := GetValueFromDict(mapping, []string{"region"}); err != nil {
		panic("Can not read region from credential file")
	} else {
		credentials.region = region.(string)
	}

	if awsKeyID, err := GetValueFromDict(mapping, []string{"aws_access_key_id"}); err != nil {
		panic("Can not read aws key from credential file")
	} else {
		credentials.awsAccessKeyID = awsKeyID.(string)
	}

	if awsSecret, err := GetValueFromDict(mapping, []string{"aws_secret_access_key"}); err != nil {
		panic("Can not read aws key from credential file")
	} else {
		credentials.awsSecretAccessKey = awsSecret.(string)
	}

	return credentials, nil
}

// CreateNewAwsClientSession creats an AWS client session
func CreateNewAwsClientSession(credentialPath string) (*session.Session, error) {
	cred, err := loadCredentialFromConfigFile(credentialPath)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cred.region),
		Credentials: credentials.NewStaticCredentials(
			cred.awsAccessKeyID, cred.awsSecretAccessKey, ""),
	})

	return sess, err
}
