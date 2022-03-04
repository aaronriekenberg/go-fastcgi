package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/utils"
)

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

type requestInfoData struct {
	RequestInfo string `json:"requestInfo"`
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

		response := &requestInfoData{
			RequestInfo: buffer.String(),
		}

		jsonText, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.CacheControlHeaderKey, utils.MaxAgeZero)

		io.Copy(w, bytes.NewReader(jsonText))
	}
}

func CreateDebugHandler(configuration *config.Configuration, serveMux *http.ServeMux) {
	serveMux.Handle("/cgi-bin/debug/request_info", requestInfoHandlerFunc())
}
