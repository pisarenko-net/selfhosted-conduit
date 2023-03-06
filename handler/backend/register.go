package backend

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pisarenko-net/selfhosted-conduit/router"
	"net/http"
)

func HandleRegisterRequests(reqRouter *mux.Router, router *router.Router) {
	reqRouter.HandleFunc("/backend/register", func(response http.ResponseWriter, request *http.Request) {
		cert := request.TLS.PeerCertificates[0]
		backendCode := router.GenerateBackendCode(cert)

		fmt.Fprintf(response, "%s", backendCode)
	}).Methods("PUT")
}
