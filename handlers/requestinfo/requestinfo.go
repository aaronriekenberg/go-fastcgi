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
	BodyContentLength int64  `json:"body_content_length"`
	Close             bool   `json:"close"`
	Host              string `json:"host"`
	Method            string `json:"method"`
	Protocol          string `json:"protocol"`
	RemoteAddress     string `json:"remote_address"`
	URL               string `json:"url"`
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
				BodyContentLength: r.ContentLength,
				Close:             r.Close,
				Host:              r.Host,
				Method:            r.Method,
				Protocol:          r.Proto,
				RemoteAddress:     r.RemoteAddr,
				URL:               urlString,
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
