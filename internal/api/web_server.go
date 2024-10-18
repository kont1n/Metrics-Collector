package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"Metrics-Collector/internal/config"
)

var err error

func Run(handler *APIHandler, sugar *zap.SugaredLogger) {
	// Запуск веб сервера
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	serverAddress := config.ParseServerConfig()
	sugar.Infoln("Web-server address: ", serverAddress)
	webSrv := StartHTTPServer(httpServerExitDone, serverAddress, handler.InitRoutes(), sugar)

	// Завершение работы веб сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = webSrv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Web server shutting down error: %s", err.Error())
	}

	httpServerExitDone.Wait()
}

func StartHTTPServer(wg *sync.WaitGroup, address string, handler http.Handler, sugar *zap.SugaredLogger) *http.Server {
	srv := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		defer wg.Done()

		if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			sugar.Fatalf("ListenAndServe() error: %v", err)
		}
	}()

	return srv
}
