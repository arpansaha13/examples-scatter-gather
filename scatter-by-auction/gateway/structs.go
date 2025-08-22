package main

import (
	"encoding/json"
	"time"
)

type ServiceResponse struct {
	ServiceID     string          `json:"service_id"`
	CorrelationID string          `json:"correlation_id"`
	Data          json.RawMessage `json:"data"`
	Timestamp     time.Time       `json:"timestamp"`
}

type AggregatedResponse struct {
	RequestID string            `json:"request_id"`
	Responses []ServiceResponse `json:"responses"`
	Count     int               `json:"count"`
	Timeout   bool              `json:"timeout"`
}

// Message published to the fanout exchange
type Message struct {
	CorrelationID string `json:"correlation_id"`
	Timestamp     int64  `json:"timestamp"`
}
