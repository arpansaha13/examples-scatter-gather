package main

type Service struct {
	Name string
	Port int
	Url  string
}

type ServiceResponse struct {
	Service   string         `json:"service"`
	Port      int            `json:"port"`
	Data      map[string]any `json:"data"`
	Timestamp string         `json:"timestamp"`
}

type AggregatedResponse struct {
	Gateway   string            `json:"gateway"`
	Port      int               `json:"port"`
	Services  []ServiceResponse `json:"services"`
	Timestamp string            `json:"timestamp"`
}
