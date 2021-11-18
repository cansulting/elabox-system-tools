package protocol

// interface for client object for connector
type ClientInterface interface {
	Join(room string) error
	IsAlive() bool
}
