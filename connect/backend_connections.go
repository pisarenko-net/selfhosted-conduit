package connect

type BackendConnections map[string]chan Message

func (c *BackendConnections) OpenBackendChannel(backendCode string) chan Message {
	connections := *c
	backendChan := make(chan Message)
	connections[backendCode] = backendChan
	return backendChan
}

func (c *BackendConnections) CloseBackendChannel(backendCode string) {
	connections := *c
	delete(connections, backendCode)
}
