package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pisarenko-net/selfhosted-conduit/connect"
	"github.com/pisarenko-net/selfhosted-conduit/util"
	"net/http"
	"time"
)

const REQUEST_TIMEOUT = 60 * time.Second

func HandleRequestRequests(reqRouter *mux.Router, connections connect.BackendConnections, responseChannels connect.ResponseChannels) {
	reqRouter.HandleFunc("/client/request/{backendCode}", func(response http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		backendCode := params["backendCode"]

		// TODO: verify if backend even exists

		// TODO: verify client is allowed to contact backend

		encryptedRequest, ok := util.GetEncryptedBodyOrFail(response, request)
		if !ok {
			return
		}

		if backendChan, ok := connections[backendCode]; ok {
			// forward client's request to the backend
			requestId, responseChan := responseChannels.OpenResponseChannel()
			defer responseChannels.CloseResponseChannel(requestId)
			clientRequestMessage := map[string]string{
				"request_id": requestId,
				"content":    encryptedRequest,
			}
			jsonBytes, _ := json.Marshal(clientRequestMessage)
			backendChan <- connect.Message{"request", string(jsonBytes)}

			// wait for the backend to respond, return to client
			clientResponseMessage := map[string]string{
				"request_id": requestId,
			}
			timeout := time.After(REQUEST_TIMEOUT)
			select {
			case responseMessage := <-responseChan:
				clientResponseMessage["content"] = responseMessage.Data
			case <-timeout:
				http.Error(response, "Backend timeout", http.StatusServiceUnavailable)
			}
			jsonBytes, _ = json.Marshal(clientResponseMessage)
			fmt.Fprintf(response, "%s\n", string(jsonBytes))
		} else {
			http.Error(response, "Backend gone", http.StatusServiceUnavailable)
			return
		}
	}).Headers("Content-Type", "text/plain").Methods("POST")
}
