package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// RegisterSQSHandler registers endpoints
func (handler *SQSHandler) RegisterSQSHandler() {
	http.HandleFunc("/server", handler.ServiceHandler)
	http.HandleFunc("/sqs", handler.SQSHandler)
	http.HandleFunc("/job", handler.HandleJob)
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

// SQSHandler handles update/get SQS URL string
func (handler *SQSHandler) SQSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		handler.QueueURL = r.URL.Query().Get("url")
		handler.ShutdownServer()
		handler.StartServer()

	} else if r.Method == "GET" {
		fmt.Fprintf(w, handler.QueueURL)
	} else {
		http.Error(w, "Not supported request method.", 405)
	}
}

func (handler *SQSHandler) HandleJob(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler.listJob(w, r)
	} else if r.Method == "POST" {
		handler.addJob(w, r)
	} else if r.Method == "DELETE" {
		handler.deleteJob(w, r)
	}
}

// addImagePatternMap add an pattern-image map
func (handler *SQSHandler) addJob(w http.ResponseWriter, r *http.Request) {
	// An example of an json body
	//
	// {
	// 	"name": "usersync",
	// 	"pattern": "s3://user_bucket/",
	// 	"image": "quay.io/cdis/fence:master",
	// 	"imageConfig": {
	//	   "indexURL": "http://index-service/"
	//	}
	// }
	//
	// Try to read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read request body; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	handler.addNewJobConfig(body)
}

func (handler *SQSHandler) deleteJob(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("pattern")
	if err := handler.deleteJobConfig(val); err != nil {
		msg := fmt.Sprintf("failed to delete a pattern; encountered error: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func (handler *SQSHandler) listJob(w http.ResponseWriter, r *http.Request) {
	str, err := handler.listAllJobConfigs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, str)
}
