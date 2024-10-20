package main

import (
	"go.uber.org/zap"
	"log"

	"Metrics-Collector/internal/api"
	"Metrics-Collector/internal/service"
	"Metrics-Collector/internal/storage"
)

var (
	err     error
	logger  *zap.Logger
	store   *storage.Store
	srv     *service.Service
	handler *api.Handler
)

func init() {
	// Подключение логирования
	logCfg := zap.NewDevelopmentConfig()
	logCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
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
