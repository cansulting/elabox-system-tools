package main

import "github.com/cansulting/elabox-system-tools/foundation/event/data"

type MyActivity struct {
}

func (instance *MyActivity) IsRunning() bool {
	return true
}

func (instance *MyActivity) OnStart(action *data.Action) error {
	
	return nil
}

func (instance *MyActivity) OnEnd() error {
	return nil
}
