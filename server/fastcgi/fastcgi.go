package fastcgi

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

	os.Remove(config.UnixSocketPath)

	listener, err := net.Listen("unix", config.UnixSocketPath)
	if err != nil {
		return nil, fmt.Errorf("net.Listen err = %w", err)
	}

	return listener, nil
}

func RunServer(
	config *config.FastCGIServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin fastcgi.RunServer UnixSocketPath = %q",
		config.UnixSocketPath,
	)

	listener, err := createListener(config)
	if err != nil {
		log.Fatalf("createListener err = %v", err)
	}

	log.Printf("before fcgi.Serve")
	error := fcgi.Serve(listener, serveHandler)
	log.Fatalf("after fcgi.Serve error = %v", error)
}
