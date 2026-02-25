package handlers

import (
	"fmt"
	"net/http"
)

// RegisterSystem
func RegisterSystem() {
	http.HandleFunc("/_status", systemStatus)
}

func systemStatus(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "Healthy"); err != nil {
		fmt.Printf("Failed to write system status response: %v\n", err)
	}
}
