package handlers

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTestJson() string {
	return `{
		"AWS": {
		  "region": "us-east-1",
		  "aws_access_key_id" : "access_key_test",
		  "aws_secret_access_key": "secret_access_key_test"
		},
		"SQS": {
		  "url": "https://sqs.us-east-1.amazonaws.com/440721843528/mySQS"
		},
		"JOBS": [
		  {
			"name": "indexing",
			"pattern": "s3://xssxs/*",
			"image": "quay.io/cdis/indexs3client:master",
			"imageConfig": {
			  "url": "http://indexd-service/",
			  "username": "test",
			  "password": "test"
			}
		  },
		  {
			"name": "usersync",
			"pattern": "s3://xssxs/user.yaml",
			"image": "quay.io/cdis/fence:master",
			"imageConfig" :{}
		  }
		]
	  }`
}

func TestUtils(t *testing.T) {

	jsonStr := makeTestJson()
	regionInf, err := GetValueFromJSON([]byte(jsonStr), []string{"AWS", "region"})
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, regionInf.(string), "us-east-1")

	accessKeyInf, err := GetValueFromJSON([]byte(jsonStr), []string{"AWS", "aws_access_key_id"})
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, accessKeyInf.(string), "access_key_test")

	_, err = GetValueFromJSON([]byte(jsonStr), []string{"AWS", "wrong_field"})
	if err == nil {
		t.Fatal(errors.New("err should not nil"))
	}

	jobInterfaces, err := GetValueFromJSON([]byte(jsonStr), []string{"JOBS"})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(jobInterfaces)
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal(b, &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	assert.Equal(t, len(jobConfigs), 2)
}

// Test that CheckIndexingJobsImageConfig does not return an error when both
// Indexd and Metadata Service creds have been configured.
func TestCheckIndexingJobsImageConfigWithIndexdAndMDSCreds(t *testing.T) {
	jobsJson :=
		`
	[
	  {
		"name": "indexing",
		"pattern": "s3://xssxs/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {
		  "url": "http://indexd-service/",
		  "username": "test",
		  "password": "test",
		  "metadataService": {
		    "url": "http://revproxy-service/mds",
		    "username": "dog",
		    "password": "paws"
		  }
		}
	  },
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.Equal(t, err, nil)
}

// Test that CheckIndexingJobsImageConfig returns an error when MDS creds
// have not been configured.
func TestCheckIndexingJobsImageConfigWithoutMDSCreds(t *testing.T) {
	jobsJson :=
		`
	[
	  {
		"name": "indexing",
		"pattern": "s3://xssxs/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {
		  "url": "http://indexd-service/",
		  "username": "test",
		  "password": "test"
		}
	  },
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.NotEqual(t, err, nil)
}

// Test that CheckIndexingJobsImageConfig returns an error when MDS creds have
// not been configured for the second indexing job.
func TestCheckIndexingJobsImageConfigWithSecondIndexingJobMissingMDSCreds(t *testing.T) {
	jobsJson :=
		`
	[
	  {
		"name": "indexing",
		"pattern": "s3://xssxs/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {
		  "url": "http://indexd-service/",
		  "username": "test",
		  "password": "test",
		  "metadataService": {
		    "url": "http://revproxy-service/mds",
		    "username": "dog",
		    "password": "paws"
		  }
		}
	  },
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  },
	  {
		"name": "indexing",
		"pattern": "s3://second-bucket/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {
		  "url": "http://indexd-service/",
		  "username": "test",
		  "password": "test"
		}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.NotEqual(t, err, nil)
}

// Test that CheckIndexingJobsImageConfig panics when the MDS password has not
// been configured.
func TestCheckIndexingJobsImageConfigWithoutMDSPassword(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expecting CheckIndexingJobsImageConfig to panic since metadataService password is missing")
		}
	}()

	jobsJson :=
		`
	[
	  {
		"name": "indexing",
		"pattern": "s3://xssxs/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {
		  "url": "http://indexd-service/",
		  "username": "test",
		  "password": "test",
		  "metadataService": {
		    "url": "http://revproxy-service/mds",
		    "username": "dog"
		  }
		}
	  },
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.NotEqual(t, err, nil)
}

// Test that CheckIndexingJobsImageConfig panics when the Indexd password has
// not been configured.
func TestCheckIndexingJobsImageConfigWithoutIndexdPassword(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expecting CheckIndexingJobsImageConfig to panic since Indexd password is missing")
		}
	}()

	jobsJson :=
		`
	[
	  {
		"name": "indexing",
		"pattern": "s3://xssxs/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {
		  "url": "http://indexd-service/",
		  "username": "test",
		  "metadataService": {
		    "url": "http://revproxy-service/mds",
		    "username": "dog",
		    "password": "paws"
		  }
		}
	  },
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.NotEqual(t, err, nil)
}

// Test that CheckIndexingJobsImageConfig panics when the indexing job's
// imageConfig is blank
func TestCheckIndexingJobsImageConfigWithoutImageConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expecting CheckIndexingJobsImageConfig to panic since indexing job imageConfig is missing")
		}
	}()

	jobsJson :=
		`
	[
	  {
		"name": "indexing",
		"pattern": "s3://xssxs/*",
		"image": "quay.io/cdis/indexs3client:master",
		"imageConfig": {}
	  },
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.NotEqual(t, err, nil)
}

// Test that CheckIndexingJobsImageConfig does not return an error when there
// is no indexing job
func TestCheckIndexingJobsImageConfigWithoutIndexingJob(t *testing.T) {
	jobsJson :=
		`
	[
	  {
		"name": "usersync",
		"pattern": "s3://xssxs/user.yaml",
		"image": "quay.io/cdis/fence:master",
		"imageConfig": {}
	  }
	]
	`
	jobConfigs := make([]JobConfig, 0)
	if err := json.Unmarshal([]byte(jobsJson), &jobConfigs); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	err := CheckIndexingJobsImageConfig(jobConfigs)
	assert.Equal(t, err, nil)
}
