package ssjdispatcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (ipm *ImagePatternMap) RegisterImagePatternMap() {
	http.HandleFunc("/pattern", ipm.handleImagePatternMap_api)
}

func (ipm *ImagePatternMap) handleImagePatternMap_api(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ipm.listImagePatternMap(w, r)
	} else if r.Method == "PUT" {
		ipm.addImagePatternMap(w, r)
	} else if r.Method == "DELETE" {
		ipm.deleteImagePatternMap(w, r)
	}

}

func (ipm *ImagePatternMap) addImagePatternMap(w http.ResponseWriter, r *http.Request) {
	// Try to read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read request body; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var mapping map[string]string

	if err := json.Unmarshal(body, &mapping); err != nil {
		msg := fmt.Sprintf("failed to read json body; encountered error: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	for pattern, quayImage := range mapping {
		ipm.AddImagePatternMap(pattern, quayImage)
	}
}

func (ipm *ImagePatternMap) deleteImagePatternMap(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("pattern")
	if err := ipm.DeleteImagePatternMap(val); err != nil {
		msg := fmt.Sprintf("failed to delete a pattern; encountered error: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}

}

func (ipm *ImagePatternMap) listImagePatternMap(w http.ResponseWriter, r *http.Request) {
	str := ""
	for pattern, handleImage := range ipm.Mapping {
		str = str + pattern + ":" + handleImage + ","
	}
	str = "{" + str + "}"
	fmt.Fprintf(w, str)
}
