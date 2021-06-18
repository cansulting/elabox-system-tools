package servicecenter

// interface for client object for connector
type ClientInterface interface {
	Join(room string) error
}
