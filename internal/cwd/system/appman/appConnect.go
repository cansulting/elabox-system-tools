package appman

import (
	"bytes"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	eventd "ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/foundation/path"
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
	location string
	//Config         data.PackageConfig
	pendingActions *eventd.ActionGroup
	client         protocol.ClientInterface
	packageId      string
	process        *os.Process
	launched       bool // true if this app was launched
}

func newAppConnect(pk *data.PackageConfig) *AppConnect {
	return &AppConnect{
		pendingActions: eventd.NewActionGroup(),
		packageId:      pk.PackageId,
		location:       path.GetAppMain(pk.PackageId, !pk.IsSystemPackage())}
}

// send pending actions
func (app *AppConnect) sendPendingActions() error {
	_, err := app.BroadcastToApp(constants.SERVICE_PENDING_ACTIONS, app.pendingActions)
	if err == nil {
		app.pendingActions.ClearAll()
	}
	return err
}

func (app *AppConnect) BroadcastToApp(action string, data interface{}) (string, error) {
	return global.Connector.BroadcastTo(app.client, action, eventd.NewAction(action, "", data))
}

func (app *AppConnect) Launch() error {
	if app.launched {
		return app.sendPendingActions()
	}
	cmd := exec.Command(app.location)
	cmd.Dir = filepath.Dir(app.location)
	app.launched = true
	go asyncRun(app, cmd)
	return nil
}

func (app *AppConnect) ForceTerminate() error {
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
	_, err := app.BroadcastToApp(constants.APP_TERMINATE, nil)
	if err != nil {
		return err
	}
	return nil
}

func asyncRun(app *AppConnect, cmd *exec.Cmd) {
	defer delete(running, app.packageId)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	if err := cmd.Start(); err != nil {
		log.Println("ERROR launching "+app.packageId, err)
		return
	}
	app.process = cmd.Process
	if err := cmd.Wait(); err != nil {
		defer log.Println("ERROR launching "+app.packageId, err)
	}
	println(app.packageId, "\n", buffer.String())
}
