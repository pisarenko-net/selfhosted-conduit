package backend

import (
	"fmt"
	"net/http"
    "github.com/pisarenko-net/selfhosted-conduit/router"
)

func HandleRegisterRequests(router *router.Router) {
    http.HandleFunc("/backend/register", func(response http.ResponseWriter, request *http.Request) {
        if request.Method != "PUT" {
            http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        cert := request.TLS.PeerCertificates[0]
        backendCode := router.GenerateBackendCode(cert)

        fmt.Fprintf(response, "%s", backendCode)
    })
}