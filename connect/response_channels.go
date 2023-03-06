package connect

import (
	"github.com/google/uuid"
)

type ResponseChannels map[string]chan Message

func (c *ResponseChannels) OpenResponseChannel() (requestId string, responseChannel chan Message) {
	requestId = uuid.New().String()
	responseChannels := *c
	responseChannel = make(chan Message)
	responseChannels[requestId] = responseChannel
	return
}

func (c *ResponseChannels) CloseResponseChannel(requestId string) string {
	responseChannels := *c
	delete(responseChannels, requestId)
	return requestId
}
