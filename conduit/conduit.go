package conduit

import (
	"crypto/tls"
	"github.com/pisarenko-net/selfhosted-conduit/handler/backend"
	"github.com/pisarenko-net/selfhosted-conduit/handler/client"
	"github.com/pisarenko-net/selfhosted-conduit/router"
	"net/http"
)

func Start(cert tls.Certificate, pathToDB string) error {
	r := router.New(pathToDB)
	defer r.Close()

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	backend.HandleRegisterRequests(r)
	backend.HandleConnectRequests(r)
	http.HandleFunc("/client/link", client.Link)

	server := &http.Server{
		Addr:      "127.0.0.1:8090",
		TLSConfig: tlsConfig,
	}

	return server.ListenAndServeTLS("", "")
}