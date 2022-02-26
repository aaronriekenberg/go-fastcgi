package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const (
	socketName            = "/var/www/run/go-fastcgi/socket"
	templatesDirectory    = "templatefiles"
	debugTemplateFile     = "debug.html"
	cacheControlHeaderKey = "cache-control"
	maxAgeZero            = "max-age=0"
)

var templates = template.Must(
	template.ParseFiles(
		filepath.Join(templatesDirectory, debugTemplateFile),
	),
)

type debugHTMLData struct {
	Title   string
	PreText string
}

func httpHeaderToString(header http.Header) string {
	var builder strings.Builder
	keys := make([]string, 0, len(header))
	for key := range header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i != 0 {
			builder.WriteRune('\n')
		}
		builder.WriteString(key)
		builder.WriteString(": ")
		fmt.Fprintf(&builder, "%v", header[key])
	}
	return builder.String()
}

func requestInfoHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer strings.Builder

		buffer.WriteString("Method: ")
		buffer.WriteString(r.Method)
		buffer.WriteRune('\n')

		buffer.WriteString("Protocol: ")
		buffer.WriteString(r.Proto)
		buffer.WriteRune('\n')

		buffer.WriteString("Host: ")
		buffer.WriteString(r.Host)
		buffer.WriteRune('\n')

		buffer.WriteString("RemoteAddr: ")
		buffer.WriteString(r.RemoteAddr)
		buffer.WriteRune('\n')

		buffer.WriteString("RequestURI: ")
		buffer.WriteString(r.RequestURI)
		buffer.WriteRune('\n')

		buffer.WriteString("URL: ")
		fmt.Fprintf(&buffer, "%#v", r.URL)
		buffer.WriteRune('\n')

		buffer.WriteString("Body.ContentLength: ")
		fmt.Fprintf(&buffer, "%v", r.ContentLength)
		buffer.WriteRune('\n')

		buffer.WriteString("Close: ")
		fmt.Fprintf(&buffer, "%v", r.Close)
		buffer.WriteRune('\n')

		buffer.WriteString("TLS: ")
		fmt.Fprintf(&buffer, "%#v", r.TLS)
		buffer.WriteString("\n\n")

		buffer.WriteString("Request Headers:\n")
		buffer.WriteString(httpHeaderToString(r.Header))

		var htmlBuilder strings.Builder
		debugHTMLData := &debugHTMLData{
			Title:   "Request Info",
			PreText: buffer.String(),
		}

		if err := templates.ExecuteTemplate(&htmlBuilder, debugTemplateFile, debugHTMLData); err != nil {
			log.Printf("error executing request info page template %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		htmlString := htmlBuilder.String()

		w.Header().Add(cacheControlHeaderKey, maxAgeZero)

		io.Copy(w, strings.NewReader(htmlString))
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("begin main socketName = %q", socketName)

	syscall.Umask(0002)
	os.Remove(socketName)

	ln, err := net.Listen("unix", socketName)
	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	serveMux := http.NewServeMux()

	serveMux.Handle("/cgi-bin/request_info", requestInfoHandlerFunc())
	//serveMux.Handle("/", http.HandlerFunc(requestInfoFunction))

	log.Printf("before fcgi.Serve")
	err = fcgi.Serve(ln, serveMux)
	log.Fatalf("after fcgi.Serve err = %v", err)
}
