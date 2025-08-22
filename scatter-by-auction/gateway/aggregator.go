package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type Aggregator struct {
	QueueName string
	Results   <-chan amqp.Delivery

	// requestChannelsMap stores buffered channels for a correlation id
	requestChannelsMap sync.Map
}

func (a *Aggregator) Start() {
	// Start consuming from results queue
	log.Println("Aggregator started, consuming from results queue")

	for msg := range a.Results {
		var response ServiceResponse
		if err := json.Unmarshal(msg.Body, &response); err != nil {
			log.Printf("Error unmarshaling response: %v", err)
			continue
		}

		logWithCorrId(
			response.CorrelationID,
			fmt.Sprintf("Received response from service %s", response.ServiceID),
		)

		// Find the corresponding channel for the correlation ID
		// And push the response to that channel
		if channelInterface, exists := a.requestChannelsMap.Load(response.CorrelationID); exists {
			if channel, ok := channelInterface.(chan ServiceResponse); ok {
				select {
				case channel <- response:
					logWithCorrId(response.CorrelationID, "Response sent to channel")
				default:
					logWithCorrId(response.CorrelationID, "Channel buffer full")
				}
			}
		} else {
			logWithCorrId(response.CorrelationID, "No channel found")
		}
	}
}
