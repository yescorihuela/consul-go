package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/status", Health)
	http.ListenAndServe(":8080", nil)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	t := struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	}{
		Status: "ok",
		Code:   200,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
