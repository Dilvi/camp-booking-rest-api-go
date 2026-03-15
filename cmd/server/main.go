package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/health", healthHandler)
	fmt.Println("server started on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string {
		"status": "ok",
	}
	w.WriteHeader(http.StatusOK)
	_ =  json.NewEncoder(w).Encode(response)
}