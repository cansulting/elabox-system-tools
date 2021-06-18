package servicecenter

import (
	"log"
	"net/http"
	"strconv"
	"time"

	socketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

////////////////////////CLASS DEFINITION///////////////
type SocketIOConnectorServer struct {
	socket *socketio.Server
	state  ConnectionType
}

func (s *SocketIOConnectorServer) GetState() ConnectionType {
	return s.state
}

// before anything else initialize this
func (s *SocketIOConnectorServer) Open() error {
	log.Println("Socket IO Server started")
	server := socketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(socketio.OnConnection, onClientConnected)
	server.On(socketio.OnDisconnection, onClientDisconnected)
	/*
		server.On("/join", func(c *socketio.Channel, channel Channel) string {
			time.Sleep(2 * time.Second)
			log.Println("Client joined to ", channel.Channel)
			return "joined to " + channel.Channel
		})*/
	s.socket = server
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	s.state = Connected
	log.Println("Starting server...")
	go http.ListenAndServe(":"+strconv.Itoa(PORT), serveMux)
	time.Sleep(time.Millisecond * 500)
	//log.Panic()

	return nil
}

func onClientConnected(socket *socketio.Channel) {
	log.Println("user connected")
	//c.Emit("/message", Message{10, "main", "using emit"})
	//c.Join("test")
	//c.BroadcastTo("test", "/message", Message{10, "main", "using broadcast"})
}

func onClientDisconnected(socket *socketio.Channel) {
	log.Println("user diiconnected")
}

// implementation for connector broadcast
func (s *SocketIOConnectorServer) Broadcast(room string, event string, data interface{}) error {
	log.Println("socketIOServer", "Broadcast", "room="+room, data)
	s.socket.BroadcastTo(room, event, data)
	return nil
}

// implementation for connector subscribe
func (s *SocketIOConnectorServer) Subscribe(room string, callback interface{}) error {
	log.Println("socketIOServer", "Subscribe", "room="+room)
	err := s.socket.On(room, callback)
	if err != nil {
		log.Panicln("socketIOServer", "Subscribe", err)
	}
	return err
}

// implementation for connector subscribe client
func (s *SocketIOConnectorServer) SubscribeClient(socket ClientInterface, room string) error {
	return socket.Join(room)
}

/// this closes the server
func (s *SocketIOConnectorServer) Close() {
	//go s.socket.()
	s.state = Disconnected
}
