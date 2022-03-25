package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
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

	for _, commandInfo := range commandConfiguration.Commands {
		apiPath := "/cgi-bin/commands/" + commandInfo.ID
		serveMux.Handle(
			apiPath,
			commandHandler.commandAPIHandlerFunc(commandInfo))
	}
}

func (commandHandler *commandHandler) getAllCommandsHandlerFunc(commandConfiguration *config.CommandConfiguration) http.HandlerFunc {
	jsonBuffer, err := json.Marshal(commandConfiguration.Commands)
	if err != nil {
		log.Fatalf("getAllCommandsHandlerFunc json.Marshal err = %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		w.Header().Add(utils.CacheControlHeaderKey, utils.CacheControlNoCache)
		io.Copy(w, bytes.NewReader(jsonBuffer))
	}
}

func (commandHandler *commandHandler) acquireCommandSemaphore(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, commandHandler.semaphoreAcquireTimeout)
	defer cancel()

	err = commandHandler.commandSemaphore.Acquire(ctx, 1)
	if err != nil {
		err = fmt.Errorf("commandHandler.acquireCommandSemaphore error calling Acquire: %w", err)
	}
	return
}

func (commandHandler *commandHandler) releaseCommandSemaphore() {
	commandHandler.commandSemaphore.Release(1)
}

type commandAPIResponse struct {
	CommandInfo     *config.CommandInfo `json:"commandInfo"`
	Now             string              `json:"now"`
	CommandDuration string              `json:"commandDuration"`
	CommandOutput   string              `json:"commandOutput"`
}

func (commandHandler *commandHandler) runCommand(ctx context.Context, commandInfo *config.CommandInfo) (response *commandAPIResponse) {
	err := commandHandler.acquireCommandSemaphore(ctx)
	if err != nil {
		response = &commandAPIResponse{
			CommandInfo:   commandInfo,
			Now:           utils.FormatTime(time.Now()),
			CommandOutput: fmt.Sprintf("%v", err),
		}
		return
	}
	defer commandHandler.releaseCommandSemaphore()

	commandStartTime := time.Now()
	rawCommandOutput, err := exec.CommandContext(
		ctx, commandInfo.Command, commandInfo.Args...).CombinedOutput()
	commandEndTime := time.Now()

	var commandOutput string
	if err != nil {
		commandOutput = fmt.Sprintf("command error %v", err)
	} else {
		commandOutput = string(rawCommandOutput)
	}

	commandDuration := fmt.Sprintf("%.9f sec",
		commandEndTime.Sub(commandStartTime).Seconds())

	response = &commandAPIResponse{
		CommandInfo:     commandInfo,
		Now:             utils.FormatTime(commandEndTime),
		CommandDuration: commandDuration,
		CommandOutput:   commandOutput,
	}
	return
}
func (commandHandler *commandHandler) commandAPIHandlerFunc(commandInfo config.CommandInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), commandHandler.requestTimeout)
		defer cancel()

		commandAPIResponse := commandHandler.runCommand(ctx, &commandInfo)

		jsonText, err := json.Marshal(commandAPIResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		w.Header().Add(utils.CacheControlHeaderKey, utils.CacheControlNoCache)
		io.Copy(w, bytes.NewReader(jsonText))
	}
}
