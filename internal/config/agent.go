package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	DefaultServer         = "localhost:8080"
	serverPath            = "/update"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2
)

var (
	flagRunAddr        string
	flagPollInterval   int
	flagReportInterval int
)

func parseAgentFlags() {
	flag.StringVar(&flagRunAddr, "a", DefaultServer, "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", DefaultReportInterval, "report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", DefaultPollInterval, "poll interval in seconds")
	flag.Parse()
}

func ParseAgentConfig() (address string, pollInterval, reportInterval time.Duration) {
	envAddress := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")

	parseAgentFlags()

	if envAddress != "" {
		address = fmt.Sprintf("http://%s%s", envAddress, serverPath)
	} else {
		address = fmt.Sprintf("http://%s%s", flagRunAddr, serverPath)
	}

	if envReportInterval != "" {
		intervalReport, err := strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			reportInterval = time.Duration(flagReportInterval) * time.Second
		}
		reportInterval = time.Duration(intervalReport) * time.Second
	} else {
		reportInterval = time.Duration(flagReportInterval) * time.Second
	}

	if envPollInterval != "" {
		intervalPoll, err := strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			pollInterval = time.Duration(flagPollInterval) * time.Second
		}
		pollInterval = time.Duration(intervalPoll) * time.Second
	} else {
		pollInterval = time.Duration(flagPollInterval) * time.Second
	}

	return address, pollInterval, reportInterval
}
