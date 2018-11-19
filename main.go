package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/uc-cdis/ssjdispatcher/ssjdispatcher"
)

func startSQSQuerying(objectHandler *ssjdispatcher.S3ObjectHandler, queueURL string) {
	for {
		time.Sleep(2 * time.Second)

		svc, err := ssjdispatcher.GetSQSSession()
		if err != nil {
			panic(err)
		}
		resp, _ := ssjdispatcher.GetSQSMessages(svc, queueURL)

		for _, message := range resp.Messages {
			// fmt.Printf("[Receive message] \n%v \n\n", message)
			err := objectHandler.HandleS3Objects(message)
			if err != nil {
				log.Println(err)
				continue
			} else {
				if ssjdispatcher.DeleteSQSMessage(svc, queueURL, message) != nil {
				} else {
					log.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
				}
			}
		}
	}
}

func main() {

	argsWithProg := os.Args
	if len(argsWithProg) < 3 {
		panic("The service requires 3. Only " + strconv.Itoa((len(argsWithProg))) + " provided")
	}

	queueURL := argsWithProg[1]
	mappingStr := argsWithProg[2]

	var mapping map[string]string
	if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
		panic(err)
	}

	// initialize S3 object handler.
	objectHandler := ssjdispatcher.S3ObjectHandler{}
	objectHandler.InitHandlerMap()
	for pattern, quayImage := range mapping {
		objectHandler.AddHandlerWithPattern(pattern, quayImage)
	}

	ssjdispatcher.RegisterSystem()
	ssjdispatcher.RegisterManager()
	go startSQSQuerying(&objectHandler, queueURL)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}
