package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/handlers"
	"github.com/aaronriekenberg/go-fastcgi/server"

	"github.com/kr/pretty"
)

func awaitShutdownSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping", s)
}
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v <config json file>", os.Args[0])
	}

	configFile := os.Args[1]

	configuration := config.ReadConfiguration(configFile)
	log.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	handler := handlers.CreateHandlers(configuration)

	server.StartServer(
		&configuration.ServerConfiguration,
		handler,
	)

	awaitShutdownSignal()
}
