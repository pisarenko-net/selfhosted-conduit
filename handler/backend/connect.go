package backend

import (
	"fmt"
	"io/ioutil"
	"net/http"
    "github.com/pisarenko-net/selfhosted-conduit/router"
    "time"
)

type testMessage struct {
    event string
    data  string
}

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
			OpenBackendChannel(response, request)
		} else {
			http.Error(response, "Provided certificate doesn't match code", http.StatusBadRequest)
			return
		}
    })
}

func OpenBackendChannel(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // TODO: this channel needs to be made available for routing
    //   client requests will make their way through here.
	var sseCh = make(chan testMessage)

	// this is just to flex the sending muscle
    go func() {
        for {
            sseCh <- testMessage{"message", "This is a message"}

            time.Sleep(3 * time.Second)
        }
    }()

    for {
        // Wait for a message to be sent to the SSE channel.
        message := <-sseCh

        // Write the SSE message to the response.
        fmt.Fprintf(w, "event: %s\n", message.event)
        fmt.Fprintf(w, "data: %s\n\n", message.data)

        // Flush the response to ensure the message is sent immediately.
        f, ok := w.(http.Flusher)
        if !ok {
            // The ResponseWriter doesn't support flushing.
            http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
            return
        }
        f.Flush()
    }
}
