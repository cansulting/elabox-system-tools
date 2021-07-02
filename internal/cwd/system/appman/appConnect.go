package appman

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	eventd "ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/internal/cwd/system/global"
	"log"
	"os"
	"os/exec"
)

/*
	struct that connects that connects and communicate to specific app
*/
type AppConnect struct {
	location string
	//Config         data.PackageConfig
	pendingActions *data.ActionGroup
	client         protocol.ClientInterface
	packageId      string
	process        *os.Process
	launched       bool // true if this app was launched
}

// send pending actions
func (app *AppConnect) sendPendingActions() error {
	_, err := app.broadcastToApp(constants.SERVICE_PENDING_ACTIONS, app.pendingActions)
	if err == nil {
		app.pendingActions.ClearAll()
	}
	return err
}

func (app *AppConnect) broadcastToApp(action string, data interface{}) (string, error) {
	return global.Connector.BroadcastTo(app.client, eventd.Action{
		Id:   action,
		Data: data,
	})
}

func (app *AppConnect) launch() error {
	if app.launched {
		return app.sendPendingActions()
	}
	app.pendingActions = eventd.NewActionGroup()
	cmd := exec.Command(app.location, "")
	app.process = cmd.Process
	app.launched = true
	go asyncRun(app.packageId, cmd)
	return nil
}

func (app *AppConnect) forceTerminate() error {
	app.launched = false
	return app.process.Kill()
}

func asyncRun(packageId string, cmd *exec.Cmd) {
	if err := cmd.Run(); err != nil {
		log.Println("ERROR ", err)
	}
	delete(running, packageId)
}
