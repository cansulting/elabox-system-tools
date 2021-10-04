package event

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"log"
	"net/http"
	"time"

	socketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

////////////////////////CLASS DEFINITION///////////////
type SocketIOServer struct {
	socket *socketio.Server
	state  data.ConnectionType
}

func (s *SocketIOServer) GetState() data.ConnectionType {
	return s.state
}

// before anything else initialize this
func (s *SocketIOServer) Open() error {
	log.Println("Socket IO Server started")
	server := socketio.NewServer(transport.GetDefaultWebsocketTransport())
	s.socket = server
	server.On(socketio.OnConnection, onClientConnected)
	server.On(socketio.OnDisconnection, onClientDisconnected)
	s.initStatus()

	http.Handle("/socket.io/", server)
	s.state = data.Connected
	//log.Println("Starting server @PORT ", constants.PORT)
	return nil
}

func onClientConnected(socket *socketio.Channel) {
	log.Println("user connected")
	//c.Emit("/message", Message{10, "main", "using emit"})
	//c.Join("test")
	//c.BroadcastTo("test", "/message", Message{10, "main", "using broadcast"})
}

func onClientDisconnected(socket *socketio.Channel) {
	log.Println("Server:onClientDisconnected", "system disconnected")
}

// implementation for connector broadcast
func (s *SocketIOServer) Broadcast(room string, event string, dataTransfer interface{}) error {
	//println("SocketIOServer", "Broadcast", "room="+room, dataTransfer)
	s.socket.BroadcastTo(room, event, dataTransfer)
	return nil
}

// implementation for connector subscribe. makes the server listen to specific room
func (s *SocketIOServer) Subscribe(room string, callback interface{}) error {
	log.Println("SocketIOServer", "Subscribe", "room="+room)
	err := s.socket.On(room, callback)
	if err != nil {
		log.Panicln("SocketIOServer", "Subscribe", err)
	}
	return err
}

// implementation for connector subscribe client
func (s *SocketIOServer) SubscribeClient(socket protocol.ClientInterface, room string) error {
	return socket.Join(room)
}

// implementation for broadcasting to specific client
func (s *SocketIOServer) BroadcastTo(client protocol.ClientInterface, method string, data interface{}) (string, error) {
	clientCast := client.(*socketio.Channel)
	return clientCast.Ack(method, data, time.Second*constants.TIMEOUT)
}

/// this closes the server
func (s *SocketIOServer) Close() error {
	return nil
}
