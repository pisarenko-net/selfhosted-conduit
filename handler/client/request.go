package client

import (
	"net/http"
    "github.com/pisarenko-net/selfhosted-conduit/connect"
    "github.com/gorilla/mux"
)

func HandleRequestRequests(reqRouter *mux.Router, connections connect.BackendConnections) {
    reqRouter.HandleFunc("/client/request/{backendCode}", func(response http.ResponseWriter, request *http.Request) {
    	params := mux.Vars(request)
        backendCode := params["backendCode"]

        // TODO: verify if backend even exists

        // TODO: verify client is allowed to contact backend

        if backendChan, ok := connections[backendCode]; ok {
        	backendChan <- connect.Message{"message", "This is a message"}
        } else {
			http.Error(response, "Backend gone", http.StatusServiceUnavailable)
            return
        }
	}).Methods("POST")
}
