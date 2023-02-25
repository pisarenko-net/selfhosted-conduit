package main

import (
	"crypto/tls"
	"github.com/pisarenko-net/selfhosted-conduit/conduit"
)

func main() {
	pathToDB := "test.db"

	certFile := "certs/server.pem"
	keyFile := "certs/server.key"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	err = conduit.Start(cert, pathToDB)
	if err != nil {
		panic(err)
	}
}
