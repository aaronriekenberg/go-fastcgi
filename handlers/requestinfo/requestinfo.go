package requestinfo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/aaronriekenberg/go-fastcgi/connection"
	"github.com/aaronriekenberg/go-fastcgi/request"
	"github.com/aaronriekenberg/go-fastcgi/utils"
)

type requestFields struct {
	ConnectionID  string `json:"connection_id"`
	RequestID     string `json:"request_id"`
	Close         bool   `json:"close"`
	ContentLength int64  `json:"content_length"`
	Host          string `json:"host"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	RemoteAddress string `json:"remote_address"`
	URL           string `json:"url"`
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

		connectionIDString := "(null)"
		connectionID, ok := r.Context().Value(connection.ConnectionIDContextKey).(connection.ConnectionID)
		if ok {
			connectionIDString = strconv.FormatInt(int64(connectionID), 10)
		}

		requestIDString := "(null)"
		requestID, ok := r.Context().Value(request.RequestIDContextKey).(request.RequestID)
		if ok {
			requestIDString = strconv.FormatInt(int64(requestID), 10)
		}

		response := &requestInfoData{
			RequestFields: requestFields{
				ConnectionID:  connectionIDString,
				RequestID:     requestIDString,
				Close:         r.Close,
				ContentLength: r.ContentLength,
				Host:          r.Host,
				Method:        r.Method,
				Protocol:      r.Proto,
				RemoteAddress: r.RemoteAddr,
				URL:           urlString,
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

func CreateRequestInfoHandler(serveMux *http.ServeMux) {
	serveMux.Handle("/cgi-bin/request_info", requestInfoHandlerFunc())
}
