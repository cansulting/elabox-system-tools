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
	isSystem bool
}

func (s *sampleact) IsRunning() bool {
	return true
}

func (s *sampleact) OnStart(action data.Action) error {
	// system update
	if s.isSystem {
		controller.RPC.CallSystem(data.NewAction(constants.SYSTEM_UPDATE_MODE, "", SAMPLE_DATA))
		return nil
	}
	if err := controller.StartActivity(data.NewAction(constants.ACTION_APP_INSTALLER, "", SAMPLE_DATA)); err != nil {
		return err
	}
	return nil
}

func (s *sampleact) OnEnd() error {
	return nil
}

// test in launching activity via broadcast
func TestNormalAppUpdate(test *testing.T) {
	var err error
	controller, err = app.NewController(&sampleact{isSystem: false}, nil)
	if err != nil {
		test.Error(err)
	}
	if err := app.RunApp(controller); err != nil {
		test.Error(err)
	}
}

// test system update
func TestSystemUpdate(test *testing.T) {
	var err error
	controller, err = app.NewController(&sampleact{isSystem: true}, nil)
	if err != nil {
		test.Error(err)
	}
	if err := app.RunApp(controller); err != nil {
		test.Error(err)
	}
}
