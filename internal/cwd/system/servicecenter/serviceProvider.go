package servicecenter

import (
	app "ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"log"
)

// function that listens to client for service center request
func OnRecievedRequest(
	client protocol.ClientInterface,
	action data.Action) string {
	switch action.Id {
	// client wants to register itself/other to a service
	case constants.SERVICE_REGISTER:
		serviceData, ok := action.Data.(*app.ServiceData)
		if ok {
			registerClientService(*serviceData)
		} else {
			log.Fatal("Invalid data")
			return "Invalid data"
		}
		break
	case constants.SERVICE_UNREGISTER:
		unregisterService(action.PackageId)
		return ""
	case constants.SERVICE_CHANGE_STATE:
		return updateServiceState(action.PackageId, action.DataToInt())
	default:

	}

	return "unknown"
}

// use to register a service
// @serviceId: the packageId or the service id
func registerClientService(
	data app.ServiceData) {
	log.Println("Requests:registerClientService", "for", data.Id)
	// TODO
}

func unregisterService(serviceId string) {
	log.Println("Requests:unregisterService", "for", serviceId)
}

func updateServiceState(serviceId string, state int) string {
	log.Println("Requests:updateServiceState", "for", serviceId, state)
	return "received"
}
