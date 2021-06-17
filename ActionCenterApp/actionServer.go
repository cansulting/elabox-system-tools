package actioncenterapp
import "ela.services/ActionCenterLib"
////////////////////////GLOBAL METHODS///////////////

var currentInstance ActionServer

func GetInstance() *ActionServer {
	return &currentInstance
}

func RunServer() {
}

////////////////////////CLASS DEFINITION///////////////
type ActionServer struct {
	connector actioncenter.ConnectorServer
}

// before anything else initialize this
func (s *ActionServer) init() {
	s.connector = actioncenter.CreateServerConnector()
}

/// this closes the server
func (s *ActionServer) Close() {
	//go s.socket.Close()
}
