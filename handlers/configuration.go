package handlers

const (
	CREDENTIAL_PATH       = "./credentials.json" // AWS credential
	JOB_NUM_MAX     int   = 5                    // number of maximum running jobs
	GRACE_PERIOD    int64 = 10                   // grace period in seconds before a job is deleted
)
