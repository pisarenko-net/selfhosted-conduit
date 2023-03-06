package client

import (
	"fmt"
	"net/http"
)

func Link(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Received link request\n")
}
