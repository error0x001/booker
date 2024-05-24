package main

import (
	"booker/internal/api/rest"
	"booker/internal/config"
	"booker/internal/domain/repos"
	"booker/internal/infrastructure/logging"
	infra "booker/internal/infrastructure/storage"
	"booker/internal/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func initRepositories(conf config.Config) (repos.OrderRepository, repos.AvailabilityRepository) {
	orderRepo, availabilityRepo, err := infra.CreateRepositories(conf)
	if err != nil {
		logging.LogPanicf(err.Error())
	}
	return orderRepo, availabilityRepo
}

func createHTTPServer(conf config.Config, orderRepo repos.OrderRepository, availabilityRepo repos.AvailabilityRepository) *http.Server {
	handler := rest.NewHandler(
		service.NewOrderService(orderRepo, availabilityRepo),
		service.NewSearchService(availabilityRepo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/availability", handler.GetAvailability)
	mux.HandleFunc("/orders", handler.CreateOrder)

	httpServer := &http.Server{ //nolint:exhaustruct
		Addr:    fmt.Sprintf(":%d", conf.ServerPort),
		Handler: mux,
	}
	return httpServer
}

func runServer(httpServer *http.Server, conf config.Config) int {
	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		logging.LogInfof("Server listening on localhost:%d", conf.ServerPort)
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			logging.LogInfof("Server closed")
			return 0
		}

		logging.LogErrorf("Server failed: %s", err)
		return 1

	case sig := <-shutdownSignals:
		logging.LogInfof("Received signal: %s, initiating shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeoutSec)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			logging.LogErrorf("Server shutdown failed: %s", err)
			return 1
		}

		logging.LogInfof("Server gracefully stopped")
		return 0
	}
}

func main() {
	conf := config.NewConfig()
	orderRepo, availabilityRepo := initRepositories(conf)
	httpServer := createHTTPServer(conf, orderRepo, availabilityRepo)
	exitCode := runServer(httpServer, conf)
	os.Exit(exitCode)
}
