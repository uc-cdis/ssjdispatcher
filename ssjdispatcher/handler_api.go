package ssjdispatcher

import (
	"fmt"
	"net/http"
)

// RegisterSQSHandler registers endpoints
func (handler *SQSHandler) RegisterSQSHandler() {
	http.HandleFunc("/handler", handler.ServiceHandler)
	http.HandleFunc("/sqs", handler.SQSHandler)
}

// ServiceHandler handles stop/start the SQS querrying service
func (handler *SQSHandler) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		val := r.URL.Query().Get("start")
		if val == "true" {
			handler.Start = true
		} else if val == "false" {
			handler.Start = false
		}

	} else {
		http.Error(w, "Not supported request method.", 405)
	}

}

// SQSHandler handles upate/get SQS URL string
func (handler *SQSHandler) SQSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		handler.QueueURL = r.URL.Query().Get("url")

	} else if r.Method == "GET" {
		fmt.Fprintf(w, handler.QueueURL)

	} else {
		http.Error(w, "Not supported request method.", 405)
	}

}
