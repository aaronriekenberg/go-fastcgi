package httpserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/connection"
)

type connWrapper struct {
	connectionID connection.ConnectionID
	net.Conn
}

func (cw *connWrapper) Close() error {
	if cw == nil {
		nilConn := (net.Conn)(nil)
		return nilConn.Close()
	}

	connectionID := cw.connectionID

	log.Printf("Removing HTTP1 connection connectionID = %v", connectionID)

	connection.ConnectionManagerInstance().RemoveConnection(connectionID)

	return cw.Conn.Close()
}

type listenerWrapper struct {
	net.Listener
}

func (lw *listenerWrapper) Accept() (net.Conn, error) {
	if lw == nil {
		nilListner := (net.Listener)(nil)
		return nilListner.Accept()
	}

	conn, err := lw.Listener.Accept()

	if err != nil {
		return conn, err
	}

	connectionID := connection.ConnectionManagerInstance().AddConnection(connection.HTTP1)

	log.Printf("Accepted HTTP1 connection connectionID = %v", connectionID)

	return &connWrapper{
		connectionID: connectionID,
		Conn:         conn,
	}, nil
}

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

	return &listenerWrapper{
		Listener: listener,
	}, nil
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
