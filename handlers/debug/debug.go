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
	HTTPMethod                string              `json:"httpMethod"`
	HTTPProtocol              string              `json:"httpProtocol"`
	Host                      string              `json:"host"`
	RemoteAddress             string              `json:"remoteAddress"`
	URL                       string              `json:"url"`
	BodyContentLength         int64               `json:"bodyContentLength"`
	Close                     bool                `json:"close"`
	SingleValueRequestHeaders map[string]string   `json:"singleValueRequestHeaders"`
	MultiValueRequestHeaders  map[string][]string `json:"multiValueRequestHeaders"`
}

func httpHeaderToSingleValueMap(headers http.Header) map[string]string {
	retVal := make(map[string]string)

	for key, value := range headers {
		if len(value) == 1 {
			retVal[key] = value[0]
		}
	}

	return retVal
}

func httpHeaderToMultiValueMap(headers http.Header) map[string][]string {
	retVal := make(map[string][]string)

	for key, value := range headers {
		if len(value) != 1 {
			retVal[key] = value
		}
	}

	return retVal
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
			HTTPMethod:                r.Method,
			HTTPProtocol:              r.Proto,
			Host:                      r.Host,
			RemoteAddress:             r.RemoteAddr,
			URL:                       urlString,
			BodyContentLength:         r.ContentLength,
			Close:                     r.Close,
			SingleValueRequestHeaders: httpHeaderToSingleValueMap(r.Header),
			MultiValueRequestHeaders:  httpHeaderToMultiValueMap(r.Header),
		}

		jsonText, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		w.Header().Add(utils.CacheControlHeaderKey, utils.CacheControlNoCache)

		io.Copy(w, bytes.NewReader(jsonText))
	}
}

func CreateDebugHandler(configuration *config.Configuration, serveMux *http.ServeMux) {
	serveMux.Handle("/cgi-bin/debug/request_info", requestInfoHandlerFunc())
}
