package config

import (
	"encoding/json"
	"log"
	"os"
)

type CommandInfo struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
}

type CommandConfiguration struct {
	MaxConcurrentCommands               int64         `json:"maxConcurrentCommands"`
	RequestTimeoutMilliseconds          int           `json:"requestTimeoutMilliseconds"`
	SemaphoreAcquireTimeoutMilliseconds int           `json:"semaphoreAcquireTimeoutMilliseconds"`
	Commands                            []CommandInfo `json:"commands"`
}

type Configuration struct {
	CommandConfiguration CommandConfiguration `json:"commandConfiguration"`
}

func ReadConfiguration(configFile string) *Configuration {
	log.Printf("reading json file %v", configFile)

	source, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading %v: %v", configFile, err)
	}

	var config Configuration
	if err = json.Unmarshal(source, &config); err != nil {
		log.Fatalf("error parsing %v: %v", configFile, err)
	}

	return &config
}
