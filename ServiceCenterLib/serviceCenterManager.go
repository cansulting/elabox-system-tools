package servicecenter

/*
	ConnectorManager.go
	Class that manages the implementation of client and server
*/

// this creates client connection
func CreateClientConnector() ConnectorClient {
	return &SocketIOConnectorClient{}
}

// this creates a server connection
func CreateServerConnector() ConnectorServer {
	return &SocketIOConnectorServer{}
}
