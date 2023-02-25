package backend

import (
	"fmt"
	"io/ioutil"
	"net/http"
    "github.com/pisarenko-net/selfhosted-conduit/router"
)

func HandleConnectRequests(router *router.Router) {
    http.HandleFunc("/backend/connect", func(response http.ResponseWriter, request *http.Request) {
        if request.Method != "POST" {
            http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(response, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		cert := request.TLS.PeerCertificates[0]

		if router.VerifyBackendCode(cert, string(body)) {
			fmt.Fprintf(response, "OK")
		} else {
			http.Error(response, "Provided certificate doesn't match code", http.StatusBadRequest)
			return
		}
    })
}
