package event

import (
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/event/socket"
)

/*
	ConnectorManager.go
	Class that manages the implementation of client and server
*/

// this creates client connection
func CreateClientConnector() protocol.ConnectorClient {
	return &socket.SocketIOClient{}
}
