package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/handlers/command"
	"github.com/aaronriekenberg/go-fastcgi/handlers/requestinfo"
)

func CreateHandlers(configuration *config.Configuration) http.Handler {
	serveMux := http.NewServeMux()

	command.CreateCommandHandler(configuration, serveMux)

	requestinfo.CreateRequestInfoHandler(configuration, serveMux)

	return serveMux
}
