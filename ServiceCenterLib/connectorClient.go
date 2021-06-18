package servicecenter

type ConnectorClient interface {
	GetState() ConnectionType
	Open() error
	Close()
	// use to send request to specific service
	SendServiceRequest(serviceId string, data ActionData) (string, error)
	// use to subscribe to specific action.
	// @callback: will be called when someone broadcasted this action
	Subscribe(subs SubscriptionData, callback interface{}) (string, error)
	Broadcast(action ActionData) error
}
