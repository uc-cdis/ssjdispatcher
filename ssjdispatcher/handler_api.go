package ssjdispatcher

import (
	"net/http"
)

// RegisterSQSHandler registers endpoints
func (handler *SQSHandler) RegisterSQSHandler() {
	http.HandleFunc("/handler", handler.ServiceHandler)
}

// ServiceHandler handles stop/start the SQS querrying service
func (handler *SQSHandler) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		val := r.URL.Query().Get("start")
		if val == "true" {
			handler.StartServer()
		} else if val == "false" {
			handler.ShutdownServer()
		}

	} else {
		http.Error(w, "Not supported request method.", 405)
	}

}
