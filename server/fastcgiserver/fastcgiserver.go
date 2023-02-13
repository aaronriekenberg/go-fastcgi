package fastcgiserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
)

func createListener(
	config *config.FastCGIServerConfiguration,
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
	config *config.FastCGIServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin fastcgiserver.Run config = %+v", *config)

	listener, err := createListener(config)
	if err != nil {
		log.Fatalf("createListener err = %v", err)
	}

	log.Printf("before fcgi.Serve")
	error := fcgi.Serve(listener, serveHandler)
	log.Fatalf("after fcgi.Serve error = %v", error)
}
