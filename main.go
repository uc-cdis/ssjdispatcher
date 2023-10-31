/*
K8S Dispatcher Service
// register SQS url
curl -X PUT http://localhost:8000/sqs?url=https://sqs.us-east-1.amazonaws.com/440721843528/mySQS
// start the service/server
curl -X PUT http://localhost:8000/server?start=true
// register new job
curl -X POST http://localhost:8000/job -d '{"Name":"TEST3", "Pattern":"s3://test/*","Image":"quay.io/cdis/indexs3client:master", "ImageConfig":{}}'
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang/glog"

	"github.com/uc-cdis/ssjdispatcher/handlers"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: ssjdispatcher -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	if err := flag.Set("logtostderr", "true"); err != nil {
		log.Fatalf("Failed to set flag: %v", err)
	}
	if err := flag.Set("stderrthreshold", "INFO"); err != nil {
		log.Fatalf("Failed to set flag: %v", err)
	}
}

func main() {
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
	jsonBytes, err := handlers.ReadFile(handlers.LookupCredFile())
	if err != nil {
		glog.Errorln("Can not read credential file!")
		os.Exit(1)
	}

	var sqsURL string
	if sqs, err := handlers.GetValueFromJSON(jsonBytes, []string{"SQS", "url"}); err != nil {
		glog.Errorln("Can not read SQS url from credential file!")
		os.Exit(1)
	} else {
		sqsURL = sqs.(string)
	}

	jobInterfaces, _ := handlers.GetValueFromJSON(jsonBytes, []string{"JOBS"})

	b, err := json.Marshal(jobInterfaces)
	if err != nil {
		glog.Info("There is no jobs configured in json credential file")
	}
	jobConfigs := make([]handlers.JobConfig, 0)
	if err := json.Unmarshal(b, &jobConfigs); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if err := handlers.CheckIndexingJobsImageConfig(jobConfigs); err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	// start an SQSHandler instance
	SQSHandler := handlers.NewSQSHandler(sqsURL)

	if err := SQSHandler.StartServer(); err != nil {
		glog.Errorf("Can not start the server. Detail %s", err)
	}

	SQSHandler.JobConfigs = jobConfigs

	SQSHandler.RegisterSQSHandler()

	handlers.RegisterJob()
	handlers.RegisterSystem()

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))

	glog.Flush()
}
