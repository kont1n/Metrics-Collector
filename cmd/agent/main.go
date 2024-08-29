package main

import (
	"flag"
	"fmt"
	"time"

	"Metrics-Collector/internal/collector"
)

const (
	DefaultServer         = "localhost:8080"
	serverPath            = "/update"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2
)

var (
	flagRunAddr        string
	flagReportInterval int
	flagPollInterval   int
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", DefaultServer, "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", DefaultReportInterval, "report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", DefaultPollInterval, "poll interval in seconds")
	flag.Parse()
}

func main() {
	parseFlags()

	pollInterval := time.Duration(flagPollInterval) * time.Second
	reportInterval := time.Duration(flagReportInterval) * time.Second
	address := fmt.Sprintf("http://%s%s", flagRunAddr, serverPath)
	fmt.Println("Agent started. Sending metrics to", address)

	agent := collector.NewAgent(address, pollInterval, reportInterval)
	agent.Run()
}
