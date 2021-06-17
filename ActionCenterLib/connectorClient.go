package actioncenter

type ConnectorClient interface {
	GetState() ConnectionType
	Open() error
	Close()
	Subscribe(subs SubscriptionData) (string, error)
	Broadcast(action ActionData) error
}
