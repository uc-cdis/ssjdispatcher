package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterJob() {
	http.HandleFunc("/job/status", status)
	http.HandleFunc("/job/list", list)
}

// status checks the status of the job given UID
func status(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not supported request method.", http.StatusMethodNotAllowed)
		return
	}
	UID := r.URL.Query().Get("UID")
	if UID != "" {
		jobHandler := NewJobHandler()
		result, errUID := jobHandler.GetJobStatusByID(UID)
		if errUID != nil {
			http.Error(w, errUID.Error(), http.StatusInternalServerError)
			return
		}

		out, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(out))
	} else {
		http.Error(w, "Missing UID argument", http.StatusMultipleChoices)
		return
	}
}

// list all the jobs
func list(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not supported request method.", http.StatusMethodNotAllowed)
		return
	}

	jobHandler := NewJobHandler()
	result := jobHandler.listJobs()

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}
