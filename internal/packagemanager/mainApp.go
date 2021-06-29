package main

import "ela/foundation/app/protocol"

type mainApp struct {
}

func (m *mainApp) OnStart() error {
	return nil
}
func (m *mainApp) IsRunning() bool {
	return true
}
func (m *mainApp) OnEnd() error {
	return nil
}

func (m *mainApp) GetService() protocol.ServiceInterface {
	return nil
}
