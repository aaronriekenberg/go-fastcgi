package debug

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aaronriekenberg/go-fastcgi/config"
	"github.com/aaronriekenberg/go-fastcgi/utils"
)

type requestInfoData struct {
	HTTPMethod        string      `json:"httpMethod"`
	HTTPProtocol      string      `json:"httpProtocol"`
	Host              string      `json:"host"`
	RemoteAddress     string      `json:"remoteAddress"`
	RequestURI        string      `json:"requestURI"`
	URL               string      `json:"url"`
	BodyContentLength int64       `json:"bodyContentLength"`
	Close             bool        `json:"close"`
	RequestHeaders    http.Header `json:"requestHeaders"`
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
			RequestURI:        r.RequestURI,
			URL:               urlString,
			BodyContentLength: r.ContentLength,
			Close:             r.Close,
			RequestHeaders:    r.Header,
		}

		jsonText, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		w.Header().Add(utils.CacheControlHeaderKey, utils.MaxAgeZero)

		io.Copy(w, bytes.NewReader(jsonText))
	}
}

func CreateDebugHandler(configuration *config.Configuration, serveMux *http.ServeMux) {
	serveMux.Handle("/cgi-bin/debug/request_info", requestInfoHandlerFunc())
}
