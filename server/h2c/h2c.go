package h2c

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"golang.org/x/net/http2"
)

func createListener(
	config *config.H2CServerConfiguration,
) (net.Listener, error) {

	os.Remove(config.UnixSocketPath)

	listener, err := net.Listen("unix", config.UnixSocketPath)
	if err != nil {
		return nil, fmt.Errorf("net.Listen err = %w", err)
	}

	return listener, nil
}

func runConnectionHandler(
	conn net.Conn,
	handler http.Handler,
	http2Server *http2.Server,
) {
	log.Printf("begin h2c.runConnectionHandler")

	defer conn.Close()

	http2Server.ServeConn(
		conn,
		&http2.ServeConnOpts{
			Handler: handler,
		})

	log.Printf("end h2c.runConnectionHandler")
}

func RunServer(
	config *config.H2CServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin h2c.RunServer UnixSocketPath = %q",
		config.UnixSocketPath,
	)

	listener, err := createListener(config)
	if err != nil {
		log.Fatalf("createListener err = %v", err)
	}

	http2Server := &http2.Server{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Accept error: %v", err)
		}

		go runConnectionHandler(conn, serveHandler, http2Server)
	}
}
