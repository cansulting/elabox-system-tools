package appman

import (
	"ela/foundation/app/data"
	"ela/foundation/constants"
	eventd "ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/foundation/path"
	"ela/internal/cwd/system/appman/nodejs"
	"ela/internal/cwd/system/global"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	nodejs         *nodejs.Nodejs
}

// create new app connect
// @pk the package config.
// @client connection from package
func newAppConnect(
	pk *data.PackageConfig,
	client protocol.ClientInterface) *AppConnect {
	// initialize node js
	var node *nodejs.Nodejs
	if pk.Nodejs {
		node = &nodejs.Nodejs{Config: pk}
	}
	return &AppConnect{
		Config:         pk,
		PendingActions: eventd.NewActionGroup(),
		PackageId:      pk.PackageId,
		Location:       path.GetAppMain(pk.PackageId, !pk.IsSystemPackage()),
		RPC:            NewRPCBridge(pk.PackageId, client, global.Connector),
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
		log.Println("Launching " + app.PackageId + " nodejs")
		go app.nodejs.Run()
	}
	// binary runnning
	if app.Config.HasMainExec() {
		log.Println("Launching " + app.PackageId + " app")
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
	log.Println("Force Terminating", app.Config.PackageId)
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
	log.Println("Terminating", app.Config.PackageId)
	if app.nodejs != nil {
		if err := app.nodejs.Stop(); err != nil {
			log.Println("AppConnect nodejs "+app.PackageId, "failed to terminate.", err.Error())
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
		log.Println("ERROR launching "+app.PackageId, err)
		return
	}
	app.process = cmd.Process
	if err := cmd.Wait(); err != nil {
		defer log.Println("ERROR launching "+app.PackageId, err)
	}
}

// callback when system has log
func (n *AppConnect) Write(data []byte) (int, error) {
	print(string(data))
	return len(data), nil
}
