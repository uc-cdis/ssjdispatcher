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
	http.HandleFunc("/jobConfig", handler.HandleJobConfig)
	http.HandleFunc("/dispatchJob", handler.HandleDispatchJob)
}

// ServiceHandler handles stop/start/status the SQS querrying service
// To start the server
//		curl -X POST http://localhost/server?start=true
// To stop the server
//		curl -X POST http://localhost/server?start=false
// To check the status
// 		curl http://localhost/server
func (handler *SQSHandler) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		val := r.URL.Query().Get("start")
		if val == "true" {
			if err := handler.StartServer(); err != nil {
				fmt.Fprintf(w, fmt.Sprintf("Can not start a server. Detail %s", err))
			} else {
				fmt.Fprintf(w, "Successfully start a server")
			}
		} else if val == "false" {
			if err := handler.ShutdownServer(); err != nil {
				fmt.Fprintf(w, fmt.Sprintf("Can not shutdown the server. Detail %s", err))
			} else {
				fmt.Fprintf(w, "Successfully shutdown the server")
			}
		}

	} else if r.Method == "GET" {
		if handler.Server != nil {
			fmt.Fprintf(w, "SQS server is running")
		} else {
			fmt.Fprintf(w, "SQS server is not running")
		}
	} else {
		http.Error(w, "Not supported request method.", 405)
	}

}

// SQSHandler handles update/get SQS URL string
// to update SQS url
//		curl -X PUT http://localhost/sqs?url=http://sqs.aws.example
// to get current sqs url
//		curl http://localhost/sqs
func (handler *SQSHandler) SQSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		handler.QueueURL = r.URL.Query().Get("url")

		if err := handler.ShutdownServer(); err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Can not shutdown the server. Detail %s", err))
		} else {
			fmt.Fprintf(w, "Successfully shutdown the server")
		}

		if err := handler.StartServer(); err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Can not start a server. Detail %s", err))
		} else {
			fmt.Fprintf(w, "Successfully start a server")
		}

	} else if r.Method == "GET" {
		fmt.Fprintf(w, handler.QueueURL)
	} else {
		http.Error(w, "Not supported request method.", 405)
	}
}

// HandleJobConfig handles job config endpoints
// to add a jobConfig
//		curl -X POST http://localhost/jobConfig -d `{"name": "usersync", "pattern": "s3://bucket/usersync.yaml", "image": "quay.io/cdis/fence:master", "imageConfig":{}}`
// to delete jobConfig
// 		curl -X DELETE http://localhost/jobConfig?pattern=s3://bucket/usersync.yaml
func (handler *SQSHandler) HandleJobConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler.listJobConfigs(w, r)
	} else if r.Method == "POST" {
		handler.addJobConfig(w, r)
	} else if r.Method == "DELETE" {
		handler.deleteJobConfig(w, r)
	}
}

// addJob adds an job config
func (handler *SQSHandler) addJobConfig(w http.ResponseWriter, r *http.Request) {
	// An example of an json body
	//
	// {
	// 	"name": "indexing",
	// 	"pattern": "s3://bucket/*",
	// 	"image": "quay.io/cdis/indexs3client:master",
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
	if err = handler.handleAddNewJobConfig(body); err != nil {
		msg := fmt.Sprintf("failed to add new job config; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Successfully add a new job!")
}

func (handler *SQSHandler) deleteJobConfig(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("pattern")
	if err := handler.handleDeleteJobConfig(val); err != nil {
		msg := fmt.Sprintf("failed to delete an job config; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Successfully delete the job")
}

func (handler *SQSHandler) listJobConfigs(w http.ResponseWriter, r *http.Request) {
	str, err := handler.handleListJobConfigs()
	if err != nil {
		msg := fmt.Sprintf("failed to list job config; encountered error: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	fmt.Fprintf(w, str)
}

// HandleDispatchJob dispatch an job
func (handler *SQSHandler) HandleDispatchJob(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handler.dispatchJob(w, r)
	}
}

// addJob adds an job config
func (handler *SQSHandler) dispatchJob(w http.ResponseWriter, r *http.Request) {
	// Try to read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read request body; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if err = handler.RetryCreateIndexingJob(body); err != nil {
		msg := fmt.Sprintf("failed to dispatch an job; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Successfully dispatch a new job!")
}
