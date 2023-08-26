package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/connection"
	"github.com/aaronriekenberg/go-fastcgi/request"
)

type connWrapper struct {
	connectionID connection.ConnectionID
	net.Conn
}

func (cw *connWrapper) Close() error {
	if cw == nil {
		return fmt.Errorf("connWrapper.Close called on nil instance")
	}

	connectionID := cw.connectionID

	slog.Info("Removing HTTP1 connection",
		slog.Any("connectionID", connectionID))

	connection.ConnectionManagerInstance().RemoveConnection(connectionID)

	return cw.Conn.Close()
}

type listenerWrapper struct {
	net.Listener
}

func (lw *listenerWrapper) Accept() (net.Conn, error) {
	if lw == nil {
		return nil, fmt.Errorf("listenerWrapper.Accept called on nil instance")
	}

	conn, err := lw.Listener.Accept()

	if err != nil {
		return conn, err
	}

	connectionID := connection.ConnectionManagerInstance().AddConnection(connection.HTTP1)

	slog.Info("Accepted HTTP1 connection",
		slog.Any("connectionID", connectionID))

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
	config config.HTTPServerConfiguration,
	handler http.Handler,
) {
	slog.Info("begin httpserver.Run",
		slog.String("config", fmt.Sprintf("%+v", config)))

	listener, err := createListener(&config)
	if err != nil {
		slog.Error("createListener error",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
	connectionManager := connection.ConnectionManagerInstance()

	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectionID, ok := r.Context().Value(connection.ConnectionIDContextKey).(connection.ConnectionID)
		if ok {
			connectionManager.IncrementRequestsForConnection(connectionID)
		}

		requestID := request.RequestIDFactoryInstance().NextRequestID()

		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, request.RequestIDContextKey, requestID))

		handler.ServeHTTP(w, r)
	})

	httpServer := &http.Server{
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			connWrapper, ok := c.(*connWrapper)
			if ok {
				connectionID := connWrapper.connectionID
				return context.WithValue(ctx, connection.ConnectionIDContextKey, connectionID)
			}
			return ctx
		},
		Handler: wrappedHandler,
	}

	httpServer.Serve(listener)
}
