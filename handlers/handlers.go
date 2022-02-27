package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/handlers/command"
	"github.com/aaronriekenberg/go-fastcgi/handlers/debug"
)

func CreateHandlers(configuration *config.Configuration) http.Handler {
	serveMux := http.NewServeMux()

	command.CreateCommandHandler(configuration, serveMux)

	debug.CreateDebugHandler(configuration, serveMux)

	return serveMux
}
