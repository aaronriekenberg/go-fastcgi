package server

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"syscall"
)

const (
	socketName = "/var/www/run/go-fastcgi/socket"
)

func RunServer(serveHandler http.Handler) error {
	log.Printf("begin RunServer socketName = %q", socketName)

	os.Remove(socketName)

	// needed so group www has rwx permission on the socket.
	syscall.Umask(0002)

	listener, err := net.Listen("unix", socketName)
	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	log.Printf("before fcgi.Serve")
	return fcgi.Serve(listener, serveHandler)
}
