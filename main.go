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

	argsWithProg := os.Args
	if len(argsWithProg) > 1 {
		queueURL = argsWithProg[1]
	}

	if len(argsWithProg) > 2 {
		mappingStr = argsWithProg[2]
	}

	SQSHandler := ssjdispatcher.NewSQSHandler(queueURL, true)

	SQSHandler.StartServer()
	defer SQSHandler.Server.Shutdown(context.Background())

	if err := SQSHandler.PatternMap.AddImagePatternMapFromJson(mappingStr); err != nil {
		log.Printf("Can not add pattern map from json %s", err)
	}

	SQSHandler.RegisterSQSHandler()
	SQSHandler.PatternMap.RegisterImagePatternMap()

	ssjdispatcher.RegisterJob()
	ssjdispatcher.RegisterSystem()

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
