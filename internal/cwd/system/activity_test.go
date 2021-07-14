package main

import (
	"ela/foundation/app"
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"testing"
)

const SAMPLE_PACKAGE = "ela.system.installer"
const SAMPLE_DATA = `C:\Users\Jhoemar\Documents\Projects\Elabox\system-tools\internal\builds\packages\packageinstaller.ela`

var controller *app.Controller

type sampleact struct {
}

func (s *sampleact) IsRunning() bool {
	return true
}

func (s *sampleact) OnStart(action data.Action) error {
	if err := controller.StartActivity(data.NewAction(
		constants.ACTION_APP_INSTALLER, SAMPLE_PACKAGE, SAMPLE_DATA)); err != nil {
		return err
	}
	return nil
}

func (s *sampleact) OnEnd() error {
	return nil
}

// test in launching activity via broadcast
func TestStartActivity(test *testing.T) {
	var err error
	controller, err = app.NewController(&sampleact{}, nil)
	if err != nil {
		test.Error(err)
	}
	if err := app.RunApp(controller); err != nil {
		test.Error(err)
	}
}
