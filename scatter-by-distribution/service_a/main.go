package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", dataHandler)

	fmt.Println("Service A starting on port 4010...")
	if err := http.ListenAndServe(":4010", nil); err != nil {
		fmt.Printf("Error starting service A: %v\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service A running on port 4010")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	dummyData := map[string]any{
		"message": "Hello from Service A",
	}

	response := ServiceResponse{
		Service:   "A",
		Port:      4010,
		Data:      dummyData,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
