package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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

	requestTimeout, err := time.ParseDuration(commandConfiguration.RequestTimeoutDuration)
	if err != nil {
		slog.Error("error parsing RequestTimeoutDuration",
			slog.String("RequestTimeoutDuration", commandConfiguration.RequestTimeoutDuration),
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	semaphoreAcquireTimeout, err := time.ParseDuration(commandConfiguration.SemaphoreAcquireTimeoutDuration)
	if err != nil {
		slog.Error("error parsing SemaphoreAcquireTimeoutDuration",
			slog.String("SemaphoreAcquireTimeoutDuration", commandConfiguration.SemaphoreAcquireTimeoutDuration),
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("CreateCommandHandler",
		slog.Group("config",
			slog.String("requestTimeout", requestTimeout.String()),
			slog.String("semaphoreAcquireTimeout", semaphoreAcquireTimeout.String())),
	)

	commandHandler := &commandHandler{
		commandSemaphore:        semaphore.NewWeighted(commandConfiguration.MaxConcurrentCommands),
		requestTimeout:          requestTimeout,
		semaphoreAcquireTimeout: semaphoreAcquireTimeout,
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
		slog.Error("getAllCommandsHandlerFunc json.Marhsal error",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
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
	CommandInfo                 *config.CommandInfo `json:"command_info"`
	Now                         string              `json:"now"`
	CommandDurationMilliseconds int64               `json:"command_duration_ms"`
	CommandOutput               string              `json:"command_output"`
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

	commandDuration := commandEndTime.Sub(commandStartTime)

	var commandOutput string
	if err != nil {
		commandOutput = fmt.Sprintf("command error %v", err)
	} else {
		commandOutput = string(rawCommandOutput)
	}

	response = &commandAPIResponse{
		CommandInfo:                 commandInfo,
		Now:                         utils.FormatTime(commandEndTime),
		CommandDurationMilliseconds: commandDuration.Milliseconds(),
		CommandOutput:               commandOutput,
	}
	return
}

func (commandHandler *commandHandler) commandAPIHandlerFunc(commandInfo config.CommandInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), commandHandler.requestTimeout)
		defer cancel()

		commandAPIResponse := commandHandler.runCommand(ctx, &commandInfo)

		jsonText, err := json.Marshal(commandAPIResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		io.Copy(w, bytes.NewReader(jsonText))
	}
}
