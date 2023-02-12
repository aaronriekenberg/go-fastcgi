package server

import (
	"net/http"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/server/fastcgi"
	"github.com/aaronriekenberg/go-fastcgi/server/h2c"
)

func StartServer(
	serverConfiguration *config.ServerConfiguration,
	serveHandler http.Handler,
) {
	switch {
	case serverConfiguration.FastCGIServerConfiguration != nil:
		go fastcgi.RunServer(serverConfiguration.FastCGIServerConfiguration, serveHandler)

	case serverConfiguration.H2CServerConfiguration != nil:
		go h2c.RunServer(serverConfiguration.H2CServerConfiguration, serveHandler)
	}
}
