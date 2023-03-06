package conduit

import (
	"crypto/tls"
	"github.com/gorilla/mux"
	"github.com/pisarenko-net/selfhosted-conduit/connect"
	"github.com/pisarenko-net/selfhosted-conduit/handler/backend"
	"github.com/pisarenko-net/selfhosted-conduit/handler/client"
	"github.com/pisarenko-net/selfhosted-conduit/router"
	"net/http"
)

func Start(cert tls.Certificate, pathToDB string) error {
	reqRouter := mux.NewRouter()

	backendRouter := router.New(pathToDB)
	defer backendRouter.Close()

	backendConnections := make(connect.BackendConnections)
	responseChannels := make(connect.ResponseChannels)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	backend.HandleRegisterRequests(reqRouter, backendRouter)
	backend.HandleConnectRequests(reqRouter, backendRouter, backendConnections)
	backend.HandleResponseRequests(reqRouter, backendRouter, responseChannels)

	client.HandleRequestRequests(reqRouter, backendConnections, responseChannels)
	http.HandleFunc("/client/link", client.Link)

	server := &http.Server{
		Addr:      "127.0.0.1:8090",
		TLSConfig: tlsConfig,
		Handler:   reqRouter,
	}

	return server.ListenAndServeTLS("", "")
}
