package main

import (
	"ela/foundation/app/protocol"
)

type appmanager struct {
}

type applicationService struct {
}

func (m *appmanager) IsRunning() bool {
	return true
}

func (m *appmanager) GetService() protocol.ServiceInterface {
	return &applicationService{}
}

func (m *appmanager) OnStart() error {
	return nil
}

func (m *appmanager) OnEnd() error {
	return nil
}

func (m *applicationService) OnStart() error {
	return nil
}
func (m *applicationService) IsRunning() bool {
	return true
}
func (m *applicationService) OnEnd() error {
	return nil
}
