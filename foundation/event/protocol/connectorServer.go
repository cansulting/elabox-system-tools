package protocol

import (
	"ela/foundation/event/data"
	"ela/foundation/system"
)

// interface for service communication to clients
type ConnectorServer interface {
	GetState() data.ConnectionType
	Open() error
	SetStatus(status system.Status, data interface{}) error
	GetStatus() string
	/// send data to all room
	Broadcast(room string, event string, data interface{}) error
	/// send service response to client
	BroadcastTo(client ClientInterface, method string, data interface{}) (string, error)
	/// server listen to room
	Subscribe(room string, callback interface{}) error
	/// make the client listen to room
	SubscribeClient(socketClient ClientInterface, room string) error
	Close() error
}
