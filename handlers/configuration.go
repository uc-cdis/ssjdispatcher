package handlers

import (
	"os"
	"strconv"
)

// LookupCredFile looks up the credential file
func LookupCredFile() string {
	val, found := os.LookupEnv("CRED_FILE")
	if !found {
		val = "/credentials.json"
	}
	return val
}

// GetMaxJobConfig returns maximum number of jobs allowed
// to run simultanously
func GetMaxJobConfig() int {
	maxJobNum, err := strconv.Atoi(os.Getenv("JOB_NUM_MAX"))
	// set to 10 if there is no env varibale
	if err != nil {
		maxJobNum = 10
	}
	return maxJobNum
}

// GetCleanupTime returns the cleanuo time for completed jobs
func GetCleanupTime() int {
	cleanupTime, err := strconv.Atoi(os.Getenv("CLEANUP_TIME"))
	if err != nil {
		cleanupTime = 120
	}
	return cleanupTime
}

const (
	GRACE_PERIOD int64 = 1 // grace period in seconds before a job is deleted
)
