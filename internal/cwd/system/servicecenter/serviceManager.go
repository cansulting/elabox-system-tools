package servicecenter

import (
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/internal/cwd/system/global"
	"log"
)

var services map[string]*ServiceConnect = make(map[string]*ServiceConnect)

// use to register a service
// @serviceId: the packageId or the service id
func onServiceOpen(
	provider protocol.ClientInterface,
	packageId string) string {
	//if err := global.Connector.SubscribeClient(client, service); err != nil {
	//	log.Fatal("Invalid data")
	//	return "Invalid data"
	//}
	services[packageId] = NewServiceConnect(packageId, provider, global.Connector)
	log.Println("serviceManager:", packageId, "is now serving")
	return ""
}

func onServiceClose(client protocol.ClientInterface, packageId string) string {
	log.Println("serviceManager:", packageId, "is now closed.")
	delete(services, packageId)
	return ""
}

func updateServiceState(client protocol.ClientInterface, action data.Action) string {
	log.Println("Requests:updateServiceState", "for")
	return "received"
}
