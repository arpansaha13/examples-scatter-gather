package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

const (
	serviceID = "service_b"
	port      = "4020"
)

func main() {
	conn, ch, err := connectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer ch.Close()
	defer conn.Close()

	log.Println("Service B connected to RabbitMQ")

	// Start consuming messages from the exchange
	go consumeMessages(ch)

	// Setup HTTP server
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", dataHandler)

	log.Printf("Service B starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service B is running on port 4020"))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Generate dummy data
	dummyData := DummyData{
		Message:   "Hello from Service B",
		ServiceID: serviceID,
		Random:    rand.Intn(1000),
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(dummyData); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func connectToRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// Declare queue for this service
	_, err = ch.QueueDeclare(
		SERVICE_B_QUEUE_NAME, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		SERVICE_B_QUEUE_NAME, // queue name
		"",                   // routing key
		FANOUT_EXCHANGE_NAME, // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}

func consumeMessages(ch *amqp.Channel) {
	// Consume messages from service queue
	msgs, err := ch.Consume(
		"service_b_queue", // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Printf("Failed to start consuming: %v", err)
		return
	}

	log.Println("Service B started consuming messages")

	for msg := range msgs {
		var message Message
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		log.Printf("Received message with correlation ID: %s", message.CorrelationID)

		// Generate dummy response data
		dummyData := DummyData{
			Message:   "Response from Service B",
			ServiceID: serviceID,
			Random:    rand.Intn(1000),
		}

		// Create response
		response := Response{
			ServiceID:     serviceID,
			CorrelationID: message.CorrelationID,
			Data:          mustMarshal(dummyData),
			Timestamp:     time.Now(),
		}

		// Publish response to results queue
		responseBody := mustMarshal(response)
		err := ch.Publish(
			"",        // exchange
			"results", // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        responseBody,
			},
		)
		if err != nil {
			log.Printf("Failed to publish response: %v", err)
		} else {
			log.Printf("Published response for correlation ID: %s", message.CorrelationID)
		}
	}
}

func mustMarshal(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to marshal: %v", err)
		return nil
	}
	return data
}
