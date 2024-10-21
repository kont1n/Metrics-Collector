package main

import (
	"log/slog"
	"os"

	"Metrics-Collector/internal/collector"
	"Metrics-Collector/internal/config"
)

var log *slog.Logger

func init() {
	// Подключение логирования
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
}

func main() {

	address, pollInterval, reportInterval := config.ParseAgentConfig(log)

	log.Info("Agent started",
		slog.String("Sending metrics to", address),
	)

	agent := collector.NewAgent(address, pollInterval, reportInterval, log)
	agent.Run()
}
