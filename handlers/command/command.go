package command

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/utils"

	"golang.org/x/sync/semaphore"
)

type commandHandler struct {
	commandSemaphore        *semaphore.Weighted
	requestTimeout          time.Duration
	semaphoreAcquireTimeout time.Duration
}

func CreateCommandHandler(configuration *config.Configuration, serveMux *http.ServeMux) {
	commandConfiguration := &configuration.CommandConfiguration
	commandHandler := &commandHandler{
		commandSemaphore:        semaphore.NewWeighted(commandConfiguration.MaxConcurrentCommands),
		requestTimeout:          time.Duration(commandConfiguration.RequestTimeoutMilliseconds) * time.Millisecond,
		semaphoreAcquireTimeout: time.Duration(commandConfiguration.SemaphoreAcquireTimeoutMilliseconds) * time.Millisecond,
	}

	serveMux.Handle(
		"/cgi-bin/commands",
		commandHandler.getAllCommandsHandlerFunc(commandConfiguration),
	)

	// for _, commandInfo := range commandConfiguration.Commands {
	// 	apiPath := "/cgi-bin/commands/" + commandInfo.ID
	// 	serveMux.Handle(
	// 		apiPath,
	// 		commandHandler.commandAPIHandlerFunc(commandInfo))
	// }
}

func (commandHandler *commandHandler) getAllCommandsHandlerFunc(commandConfiguration *config.CommandConfiguration) http.HandlerFunc {

	jsonBuffer, err := json.Marshal(commandConfiguration.Commands)
	if err != nil {
		log.Fatalf("getAllCommandsHandlerFunc json.Marshal err = %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		w.Header().Add(utils.CacheControlHeaderKey, utils.MaxAgeZero)
		io.Copy(w, bytes.NewReader(jsonBuffer))
	}
}
