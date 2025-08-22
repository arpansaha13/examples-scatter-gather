package main

import (
	"encoding/json"
	"time"
)

// Message represents a message received from the gateway
type Message struct {
	CorrelationID string `json:"correlation_id"`
	Timestamp     int64  `json:"timestamp"`
}

// Response represents a response to be sent back to the gateway
type Response struct {
	ServiceID     string          `json:"service_id"`
	CorrelationID string          `json:"correlation_id"`
	Data          json.RawMessage `json:"data"`
	Timestamp     time.Time       `json:"timestamp"`
}

type DummyData struct {
	Message   string `json:"message"`
	ServiceID string `json:"service_id"`
	Random    int    `json:"random"`
}
