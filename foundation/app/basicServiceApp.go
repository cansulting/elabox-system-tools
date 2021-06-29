package app

/*
	basicServiceApp.go
	Provides basic implementation of app with service.
*/
import (
	"ela/foundation/app/protocol"
)

type BasicServiceApp struct {
	Service protocol.ServiceInterface
}

// true if this app is running
func (m *BasicServiceApp) IsRunning() bool {
	return true
}

// return current service attach to this app
func (m *BasicServiceApp) GetService() protocol.ServiceInterface {
	return m.Service
}

// callback when this app was stated
func (m *BasicServiceApp) OnStart() error {
	return nil
}

// callback when this app ended
func (m *BasicServiceApp) OnEnd() error {
	return nil
}
