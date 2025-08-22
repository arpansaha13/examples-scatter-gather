package main

type ServiceResponse struct {
	Service   string         `json:"service"`
	Port      int            `json:"port"`
	Data      map[string]any `json:"data"`
	Timestamp string         `json:"timestamp"`
}
