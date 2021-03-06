// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package main

import (
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

const SAMPLE_PACKAGE = "ela.system.installer"
const SAMPLE_DATA = `../../builds/linux/system/ela.system.box`

var controller *app.Controller

type sampleact struct {
	isSystem bool
}

func (s *sampleact) IsRunning() bool {
	return true
}

func (instance *sampleact) OnStart() error {
	return nil
}

func (s *sampleact) OnPendingAction(action *data.Action) error {
	// system update
	if s.isSystem {
		controller.RPC.CallSystem(data.NewAction(constants.ACTION_APP_SYSTEM_INSTALL, "", SAMPLE_DATA))
		return nil
	}
	if err := controller.StartActivity(data.NewAction(constants.ACTION_APP_INSTALL, "", SAMPLE_DATA)); err != nil {
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

func TestSystemUpdateAction(test *testing.T) {
	con := event.CreateClientConnector()
	if err := con.Open(5); err != nil {
		test.Error(err)
		return
	}
	res, err := con.SendSystemRequest(constants.SYSTEM_SERVICE_ID,
		data.NewAction(constants.ACTION_START_ACTIVITY, "", data.NewAction(constants.ACTION_APP_SYSTEM_INSTALL, "", SAMPLE_DATA)))
	if err != nil {
		test.Error(err)
		return
	}
	test.Log(res)
}
