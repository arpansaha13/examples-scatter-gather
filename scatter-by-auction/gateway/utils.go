package main

import (
	"log"
)

func logWithCorrId(corrId string, msg string) {
	log.Printf("[Correlation ID: %s] %s", corrId, msg)
}
