package backend

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pisarenko-net/selfhosted-conduit/connect"
	"github.com/pisarenko-net/selfhosted-conduit/router"
	"io/ioutil"
	"net/http"
)

func HandleConnectRequests(reqRouter *mux.Router, router *router.Router, connections connect.BackendConnections) {
	reqRouter.HandleFunc("/backend/connect", func(response http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(response, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		cert := request.TLS.PeerCertificates[0]
		backendCode := string(body)

		if router.VerifyBackendCode(cert, backendCode) {
			backendChan := connections.OpenBackendChannel(backendCode)
			defer connections.CloseBackendChannel(backendCode)
			StartBackendLoop(backendChan, response, request)
		} else {
			http.Error(response, "Provided certificate doesn't match code", http.StatusBadRequest)
			return
		}
	}).Methods("POST")
}

func StartBackendLoop(backendChan chan connect.Message, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		// Wait for a message to be sent to the SSE channel.
		message := <-backendChan

		// Write the SSE message to the response.
		fmt.Fprintf(w, "event: %s\n", message.Event)
		fmt.Fprintf(w, "data: %s\n\n", message.Data)

		// Flush the response to ensure the message is sent immediately.
		f, ok := w.(http.Flusher)
		if !ok {
			// The ResponseWriter doesn't support flushing.
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		f.Flush()
	}
}
