package h2cserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/connection"
	"github.com/aaronriekenberg/go-fastcgi/request"
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
	defer conn.Close()

	connectionManager := connection.ConnectionManagerInstance()

	connectionID := connectionManager.AddConnection(connection.HTTP2)

	defer connectionManager.RemoveConnection(connectionID)

	log.Printf("begin h2cserver.runConnectionHandler connectionID = %v", connectionID)

	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectionManager.IncrementRequestsForConnection(connectionID)

		requestID := request.RequestIDFactoryInstance().NextRequestID()

		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, request.RequestIDContextKey, requestID))

		handler.ServeHTTP(w, r)
	})

	http2Server.ServeConn(
		conn,
		&http2.ServeConnOpts{
			Context: context.WithValue(context.Background(), connection.ConnectionIDContextKey, connectionID),
			Handler: wrappedHandler,
		},
	)

	log.Printf("end h2cserver.runConnectionHandler connectionID = %v", connectionID)
}

func Run(
	config config.H2CServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin h2cserver.Run config = %+v", config)

	listener, err := createListener(&config)
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
