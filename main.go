package main

import (
	"log"
	"net/http"

	"github.com/aaronriekenberg/go-fastcgi/handlers"
	"github.com/aaronriekenberg/go-fastcgi/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("begin main")

	serveMux := http.NewServeMux()

	serveMux.Handle("/cgi-bin/request_info", handlers.RequestInfoHandlerFunc())

	log.Fatalf("server.RunServer err = %v", server.RunServer(serveMux))
}
