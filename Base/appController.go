package base

import (
	"log"
	"time"

	client "ela.services/ServiceCenterLib"
)

/*
	appController.go
	Controls the main flow of app and services
*/

var connector client.ConnectorClient
var currentApp AppInterface
var appData AppData

func RunApp(app AppInterface, data AppData) error {
	appData = data
	currentApp = app
	initApp()

	for app.IsRunning() {
		time.Sleep(time.Second * 1)
	}
	return uninitApp()
}

func StartService(service ServiceInterface) error {
	service.OnStart()
	_, err := connector.SendServiceRequest(
		ACTIONCENTER_ID,
		client.ActionData{Action: ACTION_CHANGE_STATE, AppId: appData.Id, Data: 1})
	if err != nil {
		return err
	}
	return nil
}

func StopService(service ServiceInterface) {
	service.OnEnd()
}

func GetConnector() client.ConnectorClient {
	return connector
}

/////////////////////////PRIVATE//////////////////////

func initApp() error {
	// initialize connection to system
	connector = client.CreateClientConnector()
	err := connector.Open()
	if err != nil {
		return err
	}

	currentApp.OnStart()
	// initialize service
	service := currentApp.GetService()
	if service != nil {
		go StartService(service)
	}
	log.Println("App is now running")
	return nil
}

func uninitApp() error {
	service := currentApp.GetService()
	if service != nil && service.IsRunning() {
		StopService(service)
	}

	currentApp.OnEnd()
	log.Println("App exit")
	return nil
}
