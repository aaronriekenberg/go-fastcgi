package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"strconv"
	"syscall"

	"github.com/aaronriekenberg/go-fastcgi/config"
)

func createListener(
	serverConfiguration *config.ServerConfiguration,
) (net.Listener, error) {

	umask, err := strconv.ParseInt(serverConfiguration.UmaskOctal, 8, 0)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseInt err = %w", err)
	}
	umaskInt := int(umask)
	log.Printf("umaskInt = %03O", umaskInt)

	os.Remove(serverConfiguration.UnixSocketPath)

	// Umask idea from https://github.com/golang/go/issues/11822#issuecomment-123850227
	previousUmask := syscall.Umask(umaskInt)
	log.Printf("previousUmask = %03O", previousUmask)
	defer syscall.Umask(previousUmask)

	listener, err := net.Listen("unix", serverConfiguration.UnixSocketPath)
	if err != nil {
		return nil, fmt.Errorf("net.Listen err = %w", err)
	}

	return listener, nil
}

func runServer(
	serverConfiguration *config.ServerConfiguration,
	serveHandler http.Handler,
) {
	log.Printf("begin runServer UnixSocketPath = %q UmaskOctal = %q",
		serverConfiguration.UnixSocketPath,
		serverConfiguration.UmaskOctal,
	)

	listener, err := createListener(serverConfiguration)
	if err != nil {
		log.Fatalf("createListener err = %v", err)
	}

	log.Printf("before fcgi.Serve")
	error := fcgi.Serve(listener, serveHandler)
	log.Fatalf("after fcgi.Serve error = %v", error)
}

func StartServer(
	serverConfiguration *config.ServerConfiguration,
	serveHandler http.Handler,
) {
	go runServer(serverConfiguration, serveHandler)
}
