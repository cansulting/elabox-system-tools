// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// The app unit. Specifically controls the lifecycle of binary app.
// This can be a native app or nodejs app.

package appman

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	eventd "github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
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
	process        *os.Process              // main program process
	launched       bool                     // true if this app was launched
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
		node = &Nodejs{Config: pk, RestartOnFailure: true}
	}
	res := &AppConnect{
		Config:         pk,
		PendingActions: eventd.NewActionGroup(),
		PackageId:      pk.PackageId,
		Location:       pk.GetMainProgram(),
		RPC:            nil,
		nodejs:         node,
		process:        nil,
		Client:         client,
	}
	res.RPC = NewRPCBridge(pk.PackageId, res, global.Server.EventServer)
	return res
}

// use to check if this app is currently running
func (app *AppConnect) IsRunning() bool {
	if app.nodejs != nil {
		return app.nodejs.IsRunning()
	}
	return app.process != nil
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
	return app.RPC.Call(action, data)
}

// this launches both main program and node js
func (app *AppConnect) Launch() error {
	if app.launched {
		return app.sendPendingActions()
	}
	// node running
	if app.nodejs != nil && !app.nodejs.IsRunning() {
		global.Logger.Info().Msg("Launching " + app.PackageId + " nodejs")
		go func() {
			if err := app.nodejs.Run(); err != nil {
				global.Logger.Error().Err(err).Msg("Failed running node js " + app.PackageId)
			}
		}()
	}
	// binary runnning
	if app.Config.HasMainProgram() && app.process == nil {
		global.Logger.Info().Msg("Launching " + app.PackageId + " app")
		args := app.Config.ProgramArgs
		if args == nil {
			args = []string{}
		}
		cmd := exec.Command(app.Location, args...)
		cmd.Dir = filepath.Dir(app.Location)

		go asyncRun(app, cmd)
	}
	app.launched = true
	return nil
}

func (app *AppConnect) IsClientConnected() bool {
	return app.Client != nil && app.Client.IsAlive()
}

func (app *AppConnect) forceTerminate() error {
	global.Logger.Info().Caller().Msg("app will now terminate " + app.Config.PackageId)
	app.launched = false
	if app.nodejs != nil {
		if err := app.nodejs.Stop(); err != nil {
			global.Logger.Error().Err(err).Caller().Msg("AppConnect nodejs " + app.PackageId + "failed to terminate.")
		}
	}
	if app.process != nil {
		if err := app.process.Kill(); err != nil {
			return err
		}
		app.process = nil
	}
	return nil
}

func (app *AppConnect) ClearData() error {
	dir := app.Config.GetDataDir()
	if _, err := os.Stat(dir); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
		if err := os.MkdirAll(dir, perm.PUBLIC_WRITE); err != nil {
			return err
		}
	}
	return nil
}

func (app *AppConnect) Restart() error {
	if app.IsRunning() {
		if err := app.Terminate(); err != nil {
			return err
		}
	}
	return app.Launch()
}

// this terminate the app naturally
func (app *AppConnect) Terminate() error {
	if !app.IsRunning() {
		return nil
	}
	global.Logger.Info().Msg("Waiting for app " + app.Config.PackageId + " to terminate in " + strconv.Itoa(global.APP_TERMINATE_COUNTDOWN) + " seconds")
	if app.nodejs != nil {
		if err := app.nodejs.Stop(); err != nil {
			global.Logger.Error().Err(err).Caller().Msg("AppConnect nodejs " + app.PackageId + "failed to terminate.")
		}
	}
	if !app.IsClientConnected() {
		return nil
	}
	go func() {
		_, err := app.RPCCall(constants.APP_TERMINATE, nil)
		if err != nil {
			global.Logger.Error().Err(err).Caller().Msg("AppConnect " + app.PackageId + " failed sending terminate RPC.")
			app.forceTerminate()
		}
	}()

	time.Sleep(time.Second * time.Duration(global.APP_TERMINATE_COUNTDOWN))
	if app.IsRunning() {
		return app.forceTerminate()
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
	app.process = nil
}

// callback when system has log
func (n *AppConnect) Write(data []byte) (int, error) {
	print(string(data))
	return len(data), nil
}
