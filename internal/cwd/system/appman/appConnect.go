package appman

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	eventd "github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
)

/*
	struct that connects that connects and communicate to specific app
*/
type AppConnect struct {
	Location string
	//Config         data.PackageConfig
	PendingActions *eventd.ActionGroup
	StartedBy      string                   // which package who started this app
	Client         protocol.ClientInterface // socket who handles this app
	PackageId      string                   // package id of this app
	Config         *data.PackageConfig      // current package info
	process        *os.Process
	launched       bool // true if this app was launched
	RPC            *RPCBridge
	nodejs         *Nodejs
}

// create new app connect
// @pk the package config.
// @client connection from package
func newAppConnect(
	pk *data.PackageConfig,
	client protocol.ClientInterface) *AppConnect {
	// initialize node js
	var node *Nodejs
	if pk.Nodejs {
		node = &Nodejs{Config: pk}
	}
	return &AppConnect{
		Config:         pk,
		PendingActions: eventd.NewActionGroup(),
		PackageId:      pk.PackageId,
		Location:       path.GetAppMain(pk.PackageId, !pk.IsSystemPackage()),
		RPC:            NewRPCBridge(pk.PackageId, client, global.Server.EventServer),
		nodejs:         node,
	}
}

// send pending actions
func (app *AppConnect) sendPendingActions() error {
	_, err := app.RPCCall(constants.SERVICE_PENDING_ACTIONS, app.PendingActions)
	if err == nil {
		app.PendingActions.ClearAll()
	}
	return err
}

// sends action to app client
func (app *AppConnect) RPCCall(action string, data interface{}) (string, error) {
	return app.RPC.Call(action, data), nil
}

func (app *AppConnect) Launch() error {
	if app.launched {
		if app.Config.HasMainExec() {
			return app.sendPendingActions()
		}
	}
	// node running
	if app.nodejs != nil {
		global.Logger.Info().Msg("Launching " + app.PackageId + " nodejs")
		go func() {
			if err := app.nodejs.Run(); err != nil {
				global.Logger.Error().Err(err).Msg("Failed running node js " + app.PackageId)
			}
		}()
	}
	// binary runnning
	if app.Config.HasMainExec() {
		global.Logger.Info().Msg("Launching " + app.PackageId + " app")
		cmd := exec.Command(app.Location)
		cmd.Dir = filepath.Dir(app.Location)

		go asyncRun(app, cmd)
	}
	app.launched = true
	return nil
}

func (app *AppConnect) IsClientConnected() bool {
	return app.Client != nil
}

func (app *AppConnect) ForceTerminate() error {
	global.Logger.Info().Stack().Msg("Force Terminating " + app.Config.PackageId)
	app.launched = false
	if app.process != nil {
		if err := app.process.Kill(); err != nil {
			return err
		}
	}
	return nil
}

// this terminate the app naturally
func (app *AppConnect) Terminate() error {
	global.Logger.Info().Msg("Terminating " + app.Config.PackageId)
	if app.nodejs != nil {
		if err := app.nodejs.Stop(); err != nil {
			global.Logger.Error().Err(err).Caller().Msg("AppConnect nodejs " + app.PackageId + "failed to terminate.")
		}
	}
	if !app.IsClientConnected() {
		return nil
	}
	_, err := app.RPCCall(constants.APP_TERMINATE, nil)
	if err != nil {
		return err
	}
	return nil
}

func asyncRun(app *AppConnect, cmd *exec.Cmd) {
	defer delete(running, app.PackageId)
	//var buffer bytes.Buffer
	cmd.Stdout = app
	cmd.Stderr = app
	err := cmd.Start()
	if err != nil {
		global.Logger.Error().Err(err).Caller().Msg("ERROR launching " + app.PackageId)
		return
	}
	app.process = cmd.Process
	if err := cmd.Wait(); err != nil {
		global.Logger.Error().Err(err).Msg("ERROR launching " + app.PackageId)
	}
}

// callback when system has log
func (n *AppConnect) Write(data []byte) (int, error) {
	print(string(data))
	return len(data), nil
}
