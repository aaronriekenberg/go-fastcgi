package httpserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
)

func createListener(
	config *config.HTTPServerConfiguration,
) (net.Listener, error) {

	if config.Network == "unix" {
		os.Remove(config.ListenAddress)
	}

	listener, err := net.Listen(config.Network, config.ListenAddress)
	if err != nil {
		return nil, fmt.Errorf("net.Listen err = %w", err)
	}

	return listener, nil
}

func Run(
	config *config.HTTPServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin httpserver.Run config = %+v", *config)

	listener, err := createListener(config)
	if err != nil {
		log.Fatalf("createListener err = %v", err)
	}

	httpServer := &http.Server{Handler: serveHandler}

	httpServer.Serve(listener)
}
