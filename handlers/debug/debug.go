package debug

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/utils"
)

type requestHeaders struct {
	SingleValue map[string]string   `json:"singleValue"`
	MultiValue  map[string][]string `json:"multiValue"`
}

type requestInfoData struct {
	HTTPMethod        string         `json:"httpMethod"`
	HTTPProtocol      string         `json:"httpProtocol"`
	Host              string         `json:"host"`
	RemoteAddress     string         `json:"remoteAddress"`
	URL               string         `json:"url"`
	BodyContentLength int64          `json:"bodyContentLength"`
	Close             bool           `json:"close"`
	RequestHeaders    requestHeaders `json:"requestHeaders"`
}

func httpHeaderToRequestHeaders(headers http.Header) requestHeaders {

	requestHeaders := requestHeaders{
		SingleValue: make(map[string]string),
		MultiValue:  make(map[string][]string),
	}

	for key, value := range headers {
		if len(value) == 1 {
			requestHeaders.SingleValue[key] = value[0]
		} else {
			requestHeaders.MultiValue[key] = value
		}
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
			HTTPMethod:        r.Method,
			HTTPProtocol:      r.Proto,
			Host:              r.Host,
			RemoteAddress:     r.RemoteAddr,
			URL:               urlString,
			BodyContentLength: r.ContentLength,
			Close:             r.Close,
			RequestHeaders:    httpHeaderToRequestHeaders(r.Header),
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

func CreateDebugHandler(configuration *config.Configuration, serveMux *http.ServeMux) {
	serveMux.Handle("/cgi-bin/request_info", requestInfoHandlerFunc())
}
