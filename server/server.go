package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/server/h2cserver"
	"github.com/aaronriekenberg/go-fastcgi/server/httpserver"
)

func StartServer(
	serverConfiguration *config.ServerConfiguration,
	serveHandler http.Handler,
) {
	switch {
	case serverConfiguration.HTTPServerConfiguration != nil:
		go httpserver.Run(*serverConfiguration.HTTPServerConfiguration, serveHandler)

	case serverConfiguration.H2CServerConfiguration != nil:
		go h2cserver.Run(*serverConfiguration.H2CServerConfiguration, serveHandler)

	default:
		slog.Error("unable to find configured server")
		os.Exit(1)
	}
}
