package h2cserver

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

	if config.Network == "unix" {
		os.Remove(config.ListenAddress)
	}

	listener, err := net.Listen(config.Network, config.ListenAddress)
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
	log.Printf("begin h2cserver.runConnectionHandler")

	defer conn.Close()

	http2Server.ServeConn(
		conn,
		&http2.ServeConnOpts{
			Handler: handler,
		})

	log.Printf("end h2cserver.runConnectionHandler")
}

func Run(
	config *config.H2CServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin h2cserver.Run config = %+v", *config)

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
