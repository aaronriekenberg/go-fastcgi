package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type HTTPServerConfiguration struct {
	Network       string `json:"network"`
	ListenAddress string `json:"listenAddress"`
}

type H2CServerConfiguration struct {
	Network       string `json:"network"`
	ListenAddress string `json:"listenAddress"`
}

type ServerConfiguration struct {
	HTTPServerConfiguration *HTTPServerConfiguration `json:"httpServerConfiguration"`
	H2CServerConfiguration  *H2CServerConfiguration  `json:"h2cServerConfiguration"`
}

type CommandInfo struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
}

type CommandConfiguration struct {
	MaxConcurrentCommands           int64         `json:"maxConcurrentCommands"`
	RequestTimeoutDuration          string        `json:"requestTimeoutDuration"`
	SemaphoreAcquireTimeoutDuration string        `json:"semaphoreAcquireTimeoutDuration"`
	Commands                        []CommandInfo `json:"commands"`
}

type Configuration struct {
	CommandConfiguration CommandConfiguration `json:"commandConfiguration"`
	ServerConfiguration  ServerConfiguration  `json:"serverConfiguration"`
}

func ReadConfiguration(configFile string) *Configuration {
	logger := slog.Default().With("configFile", configFile)

	logger.Info("begin ReadConfiguration")

	source, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error("ReadFile error",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	var config Configuration
	if err = json.Unmarshal(source, &config); err != nil {
		logger.Error("json.Unmarshal error",
			slog.String("error", err.Error()))
	}

	logger.Info("end ReadConfiguration")

	return &config
}
