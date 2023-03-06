package backend

import (
	"github.com/gorilla/mux"
	"github.com/pisarenko-net/selfhosted-conduit/connect"
	"github.com/pisarenko-net/selfhosted-conduit/router"
	"github.com/pisarenko-net/selfhosted-conduit/util"
	"net/http"
)

func HandleResponseRequests(reqRouter *mux.Router, router *router.Router, responseChannels connect.ResponseChannels) {
	reqRouter.HandleFunc("/backend/response/{backendCode}/{requestId}", func(response http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		backendCode := params["backendCode"]
		requestId := params["requestId"]

		cert := request.TLS.PeerCertificates[0]

		if !router.VerifyBackendCode(cert, backendCode) {
			http.Error(response, "Provided certificate doesn't match code", http.StatusBadRequest)
			return
		}

		encryptedResponse, ok := util.GetEncryptedBodyOrFail(response, request)
		if !ok {
			return
		}

		if responseChan, ok := responseChannels[requestId]; ok {
			responseChan <- connect.Message{"response", encryptedResponse}
		} else {
			http.Error(response, "Client gone", http.StatusServiceUnavailable)
			return
		}
	}).Headers("Content-Type", "text/plain").Methods("POST")
}
