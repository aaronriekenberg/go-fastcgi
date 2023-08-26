package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/handlers"
	"github.com/aaronriekenberg/go-fastcgi/server"
)

func awaitShutdownSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	slog.Info("Signal received, stopping", "signal", s)
	os.Exit(0)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("begin main")

	if len(os.Args) != 2 {
		slog.Error("config file required as command line arument")
		os.Exit(1)
	}

	configFile := os.Args[1]

	configuration := config.ReadConfiguration(configFile)

	handler := handlers.CreateHandlers(configuration)

	server.StartServer(
		&configuration.ServerConfiguration,
		handler,
	)

	awaitShutdownSignal()
}
