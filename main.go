package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/uc-cdis/ssjdispatcher/ssjdispatcher"
)

func main() {
	queueURL := "https://sqs.us-east-1.amazonaws.com/440721843528/mySQS"
	mappingStr := "{\"pattern1\": \"quay.io/cdis/do_1\", \"pattern2\": \"quay.io/cdis/do_2\"}"

	argsWithProg := os.Args
	if len(argsWithProg) > 1 {
		queueURL = argsWithProg[1]
	}

	if len(argsWithProg) > 2 {
		mappingStr = argsWithProg[2]
	}

	SQSHandler := ssjdispatcher.NewSQSHandler(queueURL, true)

	var mapping map[string]string
	if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
		panic(err)
	}

	for pattern, quayImage := range mapping {
		SQSHandler.PatternMap.AddImagePatternMap(pattern, quayImage)
	}

	SQSHandler.RegisterSQSHandler()
	SQSHandler.PatternMap.RegisterImagePatternMap()

	ssjdispatcher.RegisterJob()
	ssjdispatcher.RegisterSystem()

	go ssjdispatcher.StartSQSQuerying(SQSHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}
