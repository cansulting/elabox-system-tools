package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
)

type MyService struct {
}

func (instance *MyService) IsRunning() bool {
	return true
}

func (instance *MyService) OnStart() error {
	controller.RPC.OnRecieved("testapp.action.GREETINGS", instance.onGreeted)
	return nil
}

func (instance *MyService) OnEnd() error {
	return nil

}

func (instance *MyService) onGreeted(client protocol.ClientInterface, action data.Action) string {
	return "hello"
}
