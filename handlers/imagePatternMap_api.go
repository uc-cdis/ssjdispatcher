package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RegisterImagePatternMap registers pattern-image map handler
func (ipm *ImagePatternMap) RegisterImagePatternMap() {
	http.HandleFunc("/pattern", ipm.handleImagePatternMap_api)
}

// handleImagePatternMap_api handles pattern-image map
func (ipm *ImagePatternMap) handleImagePatternMap_api(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ipm.listImagePatternMap(w, r)
	} else if r.Method == "PUT" {
		ipm.addImagePatternMap(w, r)
	} else if r.Method == "DELETE" {
		ipm.deleteImagePatternMap(w, r)
	}

}

// addImagePatternMap add an pattern-image map
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

// deleteImagePatternMap deletes pattern-image map
func (ipm *ImagePatternMap) deleteImagePatternMap(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("pattern")
	if err := ipm.DeleteImagePatternMap(val); err != nil {
		msg := fmt.Sprintf("failed to delete a pattern; encountered error: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}

}

// listImagePatternMap lists all pattern-image maps
func (ipm *ImagePatternMap) listImagePatternMap(w http.ResponseWriter, r *http.Request) {
	str := ""
	for pattern, handleImage := range ipm.Mapping {
		str = str + pattern + ":" + handleImage + ","
	}
	str = "{" + str + "}"
	fmt.Fprintf(w, str)
}
