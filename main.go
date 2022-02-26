package main

import "fmt"
import "log"
import "os"
import "syscall"
import "sort"
import "strings"
import "io"
import "net"
import "net/http"
import "net/http/fcgi"

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

func requestInfoFunction(w http.ResponseWriter, r *http.Request) {
	log.Printf("in requestInfoFunction r = %+v", r)

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

	io.Copy(w, strings.NewReader(buffer.String()))
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("hello world")

	syscall.Umask(0002)
	os.Remove("/var/www/run/go-fastcgi/socket")

	ln, err := net.Listen("unix", "/var/www/run/go-fastcgi/socket")
	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	serveMux := http.NewServeMux()

	serveMux.Handle("/cgi-bin/request_info", http.HandlerFunc(requestInfoFunction))
	//serveMux.Handle("/", http.HandlerFunc(requestInfoFunction))

	log.Printf("before fcgi.Serve")
	err = fcgi.Serve(ln, serveMux)
	log.Fatalf("after fcgi.Serve err = %v", err)
}
