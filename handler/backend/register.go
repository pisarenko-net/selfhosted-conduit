package backend

import (
	"fmt"
	"net/http"
    "github.com/pisarenko-net/selfhosted-conduit/router"
    "github.com/gorilla/mux"
)

func HandleRegisterRequests(reqRouter *mux.Router, router *router.Router) {
    reqRouter.HandleFunc("/backend/register", func(response http.ResponseWriter, request *http.Request) {
        cert := request.TLS.PeerCertificates[0]
        backendCode := router.GenerateBackendCode(cert)

        fmt.Fprintf(response, "%s", backendCode)
    }).Methods("PUT")
}