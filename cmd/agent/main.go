package main

import (
	"fmt"

	"Metrics-Collector/internal/collector"
	"Metrics-Collector/internal/config"
)

func main() {
	address, pollInterval, reportInterval := config.ParseAgentConfig()
	fmt.Println("Agent started. Sending metrics to", address)

	agent := collector.NewAgent(address, pollInterval, reportInterval)
	agent.Run()
}
