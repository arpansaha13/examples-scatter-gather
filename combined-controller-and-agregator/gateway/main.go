package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", scatterGatherHandler)

	fmt.Println("Gateway starting on port 4000...")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		fmt.Printf("Error starting gateway: %v\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Scatter Gather Gateway running on port 4000")
}

func scatterGatherHandler(w http.ResponseWriter, r *http.Request) {
	services := []Service{
		{"A", 4010, "http://service_a:4010/data"},
		{"B", 4020, "http://service_b:4020/data"},
	}

	var wg sync.WaitGroup
	responses := make([]ServiceResponse, len(services))
	errors := make([]error, len(services))

	for i, service := range services {
		wg.Add(1)

		go func(index int, svc Service) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			resp, err := client.Get(svc.Url)
			if err != nil {
				errors[index] = err
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errors[index] = err
				return
			}

			var serviceResp ServiceResponse
			if err := json.Unmarshal(body, &serviceResp); err != nil {
				errors[index] = err
				return
			}

			responses[index] = serviceResp
		}(i, service)
	}

	wg.Wait()

	var validResponses []ServiceResponse
	for i, resp := range responses {
		if errors[i] != nil {
			validResponses = append(validResponses, ServiceResponse{
				Service:   services[i].Name,
				Port:      services[i].Port,
				Data:      map[string]any{"error": errors[i].Error()},
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			validResponses = append(validResponses, resp)
		}
	}

	aggregated := AggregatedResponse{
		Gateway:   "Scatter Gather Gateway",
		Port:      4000,
		Services:  validResponses,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(aggregated)
}
