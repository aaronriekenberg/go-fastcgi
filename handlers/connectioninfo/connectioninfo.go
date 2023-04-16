package connectioninfo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/aaronriekenberg/go-fastcgi/connection"
	"github.com/aaronriekenberg/go-fastcgi/utils"
)

type connectionDTO struct {
	ID             int64  `json:"id"`
	ConnectionType string `json:"connection_type"`
	Age            string `json:"age"`
	CreationTime   string `json:"creation_time"`
	Requests       int64  `json:"requests"`
}

type connectionInfoResponse struct {
	NumConnections int              `json:"num_connections"`
	Connections    []*connectionDTO `json:"connections"`
}

func CreateConnectionInfoHandler(serveMux *http.ServeMux) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		connections := connection.ConnectionManagerInstance().Connections()

		connectionDTOs := make([]*connectionDTO, 0, len(connections))

		for _, connection := range connections {
			cdto := &connectionDTO{
				ID:             int64(connection.ID()),
				ConnectionType: connection.ConnectionType().String(),
				Age:            time.Since(connection.CreationTime()).Truncate(time.Millisecond).String(),
				CreationTime:   utils.FormatTime(connection.CreationTime()),
				Requests:       connection.Requests(),
			}

			connectionDTOs = append(connectionDTOs, cdto)
		}

		sort.Slice(connectionDTOs, func(i, j int) bool {
			return connectionDTOs[i].ID < connectionDTOs[j].ID
		})

		response := connectionInfoResponse{
			NumConnections: len(connectionDTOs),
			Connections:    connectionDTOs,
		}

		jsonText, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(utils.ContentTypeHeaderKey, utils.ContentTypeApplicationJSON)
		io.Copy(w, bytes.NewReader(jsonText))
	}

	serveMux.HandleFunc("/cgi-bin/connection_info", http.HandlerFunc(handler))
}
