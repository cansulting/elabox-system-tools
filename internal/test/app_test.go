package main

import (
	"os/exec"
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

const pkgid = "ela.testapp"
const samplesrc = "../../sample/simple"
const builddir = "../../sample/simple/build"
const sampleapp_bin = "../../sample/simple/build/bin/" + pkgid
const GREET_ACTION = "testapp.action.GREETINGS"

// test app packaging. Make sure ela packager is installed
func TestAppPackaging(t *testing.T) {
	// build
	cmd := exec.Command(
		"go", "build",
		"-ldflags", "-w -s",
		"-tags", "RELEASE",
		"-o", sampleapp_bin,
		samplesrc)
	output, err := cmd.CombinedOutput()
	println(string(output[:]))
	if err != nil {
		t.Error("Failed building sample app", err)
		return
	}

	// package
	cmd = exec.Command("packager", builddir+"/packager.json")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Error("Failed packaging sample app", err, string(output[:]))
		return
	}
}

// test app install. Make sure ela installer is installed
func TestAppInstall(t *testing.T) {
	cmd := exec.Command("packageinstaller", builddir+"/"+pkgid+".box")
	output, err := cmd.CombinedOutput()
	println(string(output[:]))
	if err != nil {
		t.Error("Failed installing sample app", err)
		return
	}
}

// test running service
func TestLaunchService(t *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		t.Error("Failed initializing RPC", err)
		return
	}
	res, err := handler.CallSystem(data.NewAction(constants.ACTION_START_SERVICE, pkgid, nil))
	if err != nil {
		t.Error("Failed RPC", err)
		return
	}
	resp, err := res.ToSimpleResponse()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.Code == 200 {
		return
	}
	t.Error(resp.Message)
}

func TestServiceRPC(t *testing.T) {
	handler, err := rpc.NewRPCHandlerDefault()
	if err != nil {
		t.Error("Failed initializing RPC", err)
		return
	}
	res, err := handler.CallRPC(pkgid, data.NewAction(GREET_ACTION, "", nil))
	if err != nil {
		t.Error("Failed RPC", err)
		return
	}
	println(res.ToString())
}
