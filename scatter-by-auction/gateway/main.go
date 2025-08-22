package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	aggr *Aggregator
	sq   *ScatterQueue
)

func main() {
	sq = &ScatterQueue{
		User:         RABBITMQ_USER,
		Pass:         RABBITMQ_PASS,
		Port:         RABBITMQ_PORT,
		Exchange:     FANOUT_EXCHANGE_NAME,
		ResultsQueue: RESULTS_QUEUE_NAME,
	}

	err := sq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer sq.Close()

	log.Println("Connected to RabbitMQ")

	results, err := sq.GetResultsDelivery()
	if err != nil {
		log.Fatalf("Failed to get results delivery: %v", err)
	}

	aggr = &Aggregator{QueueName: RESULTS_QUEUE_NAME, Results: results}
	go aggr.Start()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", dataHandler)

	log.Printf("Gateway starting on port %d", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Gateway running on port %d", PORT)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	correlationID := uuid.New().String()

	responseChannel := make(chan ServiceResponse, 10)
	aggr.requestChannelsMap.Store(correlationID, responseChannel)

	defer func() {
		aggr.requestChannelsMap.Delete(correlationID)
		close(responseChannel)
	}()

	message := Message{
		CorrelationID: correlationID,
		Timestamp:     time.Now().Unix(),
	}
	if err := sq.Publish(&message); err != nil {
		log.Printf("Failed to publish message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logWithCorrId(correlationID, "Published message")

	responses, timeout := waitForResponses(
		responseChannel,
		MIN_RESPONSE_COUNT,
		time.Duration(RESPONSE_TIMEOUT_SECONDS)*time.Second,
	)

	aggregatedResponse := AggregatedResponse{
		RequestID: correlationID,
		Responses: responses,
		Count:     len(responses),
		Timeout:   timeout,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(aggregatedResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logWithCorrId(
		correlationID,
		fmt.Sprintf("Request completed with %d responses, timeout: %v", len(responses), timeout),
	)
}

func waitForResponses(channel chan ServiceResponse, minResponses int, timeout time.Duration) ([]ServiceResponse, bool) {
	var responses []ServiceResponse
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case response := <-channel:
			responses = append(responses, response)
			if len(responses) >= minResponses {
				return responses, false
			}
		case <-timer.C:
			return responses, true
		}
	}
}
