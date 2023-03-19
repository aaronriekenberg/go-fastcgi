package requestinfo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/utils"
)

type requestFields struct {
	HTTPMethod        string `json:"http_method"`
	HTTPProtocol      string `json:"http_protocol"`
	Host              string `json:"host"`
	RemoteAddress     string `json:"remote_address"`
	URL               string `json:"url"`
	BodyContentLength int64  `json:"body_content_length"`
	Close             bool   `json:"close"`
}

type requestInfoData struct {
	RequestFields  requestFields     `json:"request_fields"`
	RequestHeaders map[string]string `json:"request_headers"`
}

func httpHeaderToRequestHeaders(headers http.Header) map[string]string {

	requestHeaders := make(map[string]string)

	for key, value := range headers {
		requestHeaders[key] = strings.Join(value, "; ")
	}

	return requestHeaders
}

func requestInfoHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var urlString string
		if r.URL != nil {
			urlString = r.URL.String()
		} else {
			urlString = "(null)"
		}

		response := &requestInfoData{
			RequestFields: requestFields{
				HTTPMethod:        r.Method,
				HTTPProtocol:      r.Proto,
				Host:              r.Host,
				RemoteAddress:     r.RemoteAddr,
				URL:               urlString,
				BodyContentLength: r.ContentLength,
				Close:             r.Close,
			},
			RequestHeaders: httpHeaderToRequestHeaders(r.Header),
		}

		jsonText, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)

		io.Copy(w, bytes.NewReader(jsonText))
	}
}

func CreateRequestInfoHandler(configuration *config.Configuration, serveMux *http.ServeMux) {
	serveMux.Handle("/cgi-bin/request_info", requestInfoHandlerFunc())
}
