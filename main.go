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
	"context"
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
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
}

func main() {
	jsonBytes, err := handlers.ReadFile(handlers.LookupCredFile())
	if err != nil {
		glog.Errorln("Can not read credential file!")
		return
	}

	var sqsURL string
	if sqs, err := handlers.GetValueFromJSON(jsonBytes, []string{"SQS", "url"}); err != nil {
		glog.Errorln("Can not read SQS url from credential file!")
		return
	} else {
		sqsURL = sqs.(string)
	}

	jobInterfaces, _ := handlers.GetValueFromJSON(jsonBytes, []string{"JOBS"})

	b, err := json.Marshal(jobInterfaces)
	if err != nil {
		glog.Info("There is no jobs configured in json credential file")
	}
	jobConfigs := make([]handlers.JobConfig, 0)
	json.Unmarshal(b, &jobConfigs)

	// start an SQSHandler instance
	SQSHandler := handlers.NewSQSHandler(sqsURL)
	if err := SQSHandler.StartServer(); err != nil {
		glog.Errorf("Can not start the server. Detail %s", err)
	}
	defer SQSHandler.Server.Shutdown(context.Background())

	SQSHandler.JobConfigs = jobConfigs

	SQSHandler.RegisterSQSHandler()

	handlers.RegisterJob()
	handlers.RegisterSystem()

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))

	glog.Flush()
}
