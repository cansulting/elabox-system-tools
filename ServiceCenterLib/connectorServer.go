package servicecenter

type ConnectionType int

const (
	Disconnected = iota
	Connected
	Connecting
	Disconnecting
)

// interface for service communication to clients
type ConnectorServer interface {
	GetState() ConnectionType
	Open() error
	/// send data to all room
	Broadcast(room string, event string, data interface{}) error
	/// server listen to room
	Subscribe(room string, callback interface{}) error
	/// make the client listen to room
	SubscribeClient(socketClient ClientInterface, room string) error
	Close()
}
