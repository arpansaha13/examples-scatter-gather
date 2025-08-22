package main

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type ScatterQueue struct {
	ch           *amqp.Channel
	conn         *amqp.Connection
	Exchange     string
	ResultsQueue string
	User         string
	Pass         string
	Port         int
}

func (sq *ScatterQueue) Connect() error {
	url := fmt.Sprintf("amqp://%s:%s@rabbitmq:%d/", sq.User, sq.Pass, sq.Port)

	var err error
	sq.conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	sq.ch, err = sq.conn.Channel()
	if err != nil {
		sq.conn.Close()
		return err
	}

	err = sq.declareFanoutExchange()
	if err != nil {
		sq.Close()
		return err
	}

	err = sq.declareResultsQueue()
	if err != nil {
		sq.Close()
		return err
	}

	return nil
}

func (sq *ScatterQueue) declareFanoutExchange() error {
	err := sq.ch.ExchangeDeclare(
		sq.Exchange, // name
		"fanout",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func (sq *ScatterQueue) declareResultsQueue() error {
	_, err := sq.ch.QueueDeclare(
		sq.ResultsQueue, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func (sq *ScatterQueue) Publish(message *Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return sq.ch.Publish(
		sq.Exchange, // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (sq *ScatterQueue) GetResultsDelivery() (<-chan amqp.Delivery, error) {
	results, err := sq.ch.Consume(
		sq.ResultsQueue, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (sq *ScatterQueue) Close() {
	if sq.ch != nil {
		sq.ch.Close()
	}
	if sq.conn != nil {
		sq.conn.Close()
	}
}
