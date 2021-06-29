package app

import (
	"log"
	"time"

	"ela/foundation/app/data"
	"ela/foundation/app/protocol"
	"ela/foundation/constants"
	"ela/foundation/event"
	eventdata "ela/foundation/event/data"
	eventprot "ela/foundation/event/protocol"
)

/*
	controller.go
	Controls the main flow of app and services
*/

var connector eventprot.ConnectorClient
var currentApp protocol.AppInterface
var appData data.AppData

func RunApp(app protocol.AppInterface, data data.AppData) error {
	appData = data
	currentApp = app
	initApp()

	for app.IsRunning() {
		time.Sleep(time.Second * 1)
	}
	return uninitApp()
}

// stops a service. sends awake state to service center
func StartService(service protocol.ServiceInterface) error {
	_, err := sendServiceCenterRequest(
		eventdata.Action{
			Action: constants.SERVICE_CHANGE_STATE,
			AppId:  appData.Id,
			Data:   constants.SERVICE_AWAKE})
	if err != nil {
		return err
	}
	service.OnStart()
	return nil
}

// stops a service. send a sleep state to service center
func StopService(service protocol.ServiceInterface) error {
	_, err := sendServiceCenterRequest(
		eventdata.Action{
			Action: constants.SERVICE_CHANGE_STATE,
			AppId:  appData.Id,
			Data:   constants.SERVICE_SLEEP})

	if err != nil {
		return err
	}
	service.OnEnd()
	return nil
}

func GetConnector() eventprot.ConnectorClient {
	return connector
}

/////////////////////////PRIVATE//////////////////////

// initialization process for the app
func initApp() error {
	// step: initialize connection to system
	connector = event.CreateClientConnector()
	err := connector.Open()
	if err != nil {
		return err
	}

	// start the app
	currentApp.OnStart()

	// initialize service
	service := currentApp.GetService()
	if service != nil {
		go StartService(service)
	}
	log.Println("App is now running")
	return nil
}

// uninitialize the app
func uninitApp() error {
	// stops the service
	service := currentApp.GetService()
	if service != nil && service.IsRunning() {
		StopService(service)
	}

	// end the app
	currentApp.OnEnd()
	log.Println("App exit")
	return nil
}

// this sends action to service center
func sendServiceCenterRequest(action eventdata.Action) (interface{}, error) {
	return connector.SendServiceRequest(
		constants.SERVICE_CENTER_ID,
		action)
}
