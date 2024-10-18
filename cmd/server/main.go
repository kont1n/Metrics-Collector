package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"Metrics-Collector/internal/api"
	"Metrics-Collector/internal/service"
	"Metrics-Collector/internal/storage"
)

var (
	err      error
	logLevel int
	logger   *zap.Logger
	store    *storage.Store
	srv      *service.Service
	handler  *api.ApiHandler
)

func init() {
	// Подключение к файлу переменных окружения
	if err = godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	logLevel, err = strconv.Atoi(os.Getenv("LOGGER_LEVEL"))
	if err != nil {
		log.Fatal("LOGGER_LEVEL is not set")
	}

	// Подключение логирования
	var level zapcore.Level
	level = zapcore.Level(logLevel)
	logCfg := zap.NewDevelopmentConfig()
	logCfg.Level = zap.NewAtomicLevelAt(level)
	logger, err = logCfg.Build()
	if err != nil {
		log.Fatal("Failed to build logger:", err.Error())
		return
	}
	defer logger.Sync()

	logger.Info("Initializing success")
}

func main() {
	sugar := logger.Sugar()

	store = storage.NewStore(sugar)
	srv = service.NewService(store, sugar)
	handler = api.NewHandler(srv, sugar)

	sugar.Infof("Application started")
	api.Run(handler, sugar)
	sugar.Infof("Application shutdown")
}
