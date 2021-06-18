package main

import (
	"log"

	base "ela.services/Base"
	lib "ela.services/ServiceCenterLib"
)

// function that listens to client Requests for service center request
func OnRecievedRequest(
	client lib.ClientInterface,
	data lib.ActionData) string {
	switch data.Action {
	case base.ACTION_REGISTER:
		serviceData, ok := data.Data.(*lib.ServiceData)
		if ok {
			registerClientService(*serviceData)
		} else {
			log.Fatal("Invalid data")
			return "Invalid data"
		}
		break
	case base.ACTION_UNREGISTER:
		unregisterService(data.AppId)
		return ""
	case base.ACTION_CHANGE_STATE:
		return updateServiceState(data.AppId, data.DataToInt())
	}

	return "unknown"
}

// use to register a service
// @serviceId: the packageId or the service id
func registerClientService(
	data lib.ServiceData) {
	log.Println("Service was registered ", data.Id)
	// TODO
}

func unregisterService(serviceId string) {
	log.Println("Service was unregistered ", serviceId)
}

func updateServiceState(serviceId string, state int) string {
	log.Println("Service state was updated ", serviceId, state)
	return "received"
}
