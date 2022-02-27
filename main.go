package main

import (
	"log"
	"os"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/handlers"
	"github.com/aaronriekenberg/go-fastcgi/server"

	"github.com/kr/pretty"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v <config json file>", os.Args[0])
	}

	configFile := os.Args[1]

	configuration := config.ReadConfiguration(configFile)
	log.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	handler := handlers.CreateHandlers(configuration)

	log.Fatalf("server.RunServer err = %v", server.RunServer(handler))
}
