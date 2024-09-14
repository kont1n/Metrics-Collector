package main

import (
	"fmt"
	"time"

	"Metrics-Collector/internal/collector"
)

const (
	serverHost     = "http://localhost"
	serverPort     = 8080
	serverPath     = "/update"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	address := fmt.Sprintf("%s:%d%s", serverHost, serverPort, serverPath)
	fmt.Println("Agent started. Sending metrics to", address)

	agent := collector.NewAgent(address, pollInterval, reportInterval)
	agent.Run()
}
