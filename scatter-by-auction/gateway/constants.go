package main

const PORT int16 = 4000

const (
	RABBITMQ_PORT int    = 5672
	RABBITMQ_USER string = "guest"
	RABBITMQ_PASS string = "guest"
)

const (
	RESULTS_QUEUE_NAME   string = "results"
	FANOUT_EXCHANGE_NAME string = "scatter-gather"
)

const (
	MIN_RESPONSE_COUNT       int   = 2
	RESPONSE_TIMEOUT_SECONDS int16 = 10
)
