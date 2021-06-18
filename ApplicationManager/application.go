package main

import (
	base "ela.services/Base"
)

type ApplicationManager struct {
}

type applicationService struct {
}

func (m *ApplicationManager) IsRunning() bool {
	return true
}

func (m *ApplicationManager) GetService() base.ServiceInterface {
	return &applicationService{}
}

func (m *ApplicationManager) OnStart() error {
	return nil
}

func (m *ApplicationManager) OnEnd() error {
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
