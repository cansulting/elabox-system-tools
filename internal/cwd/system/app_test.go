package main

import (
	"os"
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

const TEST_APP = "ela.logs"
const TEST_DATA = "/var/ela/data/" + TEST_APP + "/test.dat"

func TestAppRestart(test *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		test.Error(err)
		return
	}
	res, err := handler.CallSystem(data.NewAction(constants.ACTION_APP_RESTART, TEST_APP, nil))
	if err != nil {
		test.Error(err)
		return
	}
	test.Log(res)
}

func TestAppOff(test *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		test.Error(err)
		return
	}
	res, err := handler.CallSystem(data.NewAction(constants.ACTION_APP_OFF, TEST_APP, nil))
	if err != nil {
		test.Error(err)
		return
	}
	test.Log(res)
}
func TestAppOn(test *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		test.Error(err)
		return
	}
	res, err := handler.CallSystem(data.NewAction(constants.ACTION_APP_ON, TEST_APP, nil))
	if err != nil {
		test.Error(err)
		return
	}
	test.Log(res)
}
func TestAppClearData(test *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		test.Error(err)
		return
	}
	if err := os.WriteFile(TEST_DATA, []byte("test"), 0644); err != nil {
		test.Error(err)
		return
	}
	res, err := handler.CallSystem(data.NewAction(constants.ACTION_APP_CLEAR_DATA, TEST_APP, nil))
	if err != nil {
		test.Error(err)
		return
	}
	if _, err := os.Stat(TEST_DATA); err == nil {
		test.Error("data not cleared")
		return
	}
	test.Log(res)
}

func TestGetDeviceSerial(test *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		test.Error(err)
		return
	}
	res, err := handler.CallSystem(data.NewAction(constants.ACTION_APP_DEVICE_SERIAL, "", nil))
	if err != nil {
		test.Error(err)
		return
	}
	test.Log(res)
}
