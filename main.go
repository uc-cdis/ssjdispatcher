package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/uc-cdis/ssjdispatcher/ssjdispatcher"
)

func main() {
	var queueURL string
	var mappingStr string
	queueURL = "https://sqs.us-east-1.amazonaws.com/440721843528/mySQS"
	mappingStr = "{\"s3://xssxs*\": \"quay.io/cdis/simu_demo:latest\"}"

	argsWithProg := os.Args
	if len(argsWithProg) > 1 {
		queueURL = argsWithProg[1]
	}

	if len(argsWithProg) > 2 {
		mappingStr = argsWithProg[2]
	}

	SQSHandler := ssjdispatcher.NewSQSHandler(queueURL, true)
	if err := SQSHandler.PatternMap.AddImagePatternMapFromJson(mappingStr); err != nil {
		log.Printf("Can not add pattern map from json %s", err)
	}

	SQSHandler.StartServer()
	defer SQSHandler.Server.Shutdown(context.Background())

	SQSHandler.RegisterSQSHandler()
	SQSHandler.PatternMap.RegisterImagePatternMap()

	ssjdispatcher.RegisterJob()
	ssjdispatcher.RegisterSystem()

	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}
