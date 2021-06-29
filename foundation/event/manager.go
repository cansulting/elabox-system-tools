package event

import (
	"ela/foundation/event/protocol"
	"ela/foundation/event/socket"
)

/*
	ConnectorManager.go
	Class that manages the implementation of client and server
*/

// this creates client connection
func CreateClientConnector() protocol.ConnectorClient {
	return &socket.SocketIOClient{}
}

// this creates a server connection
func CreateServerConnector() protocol.ConnectorServer {
	return &socket.SocketIOServer{}
}
