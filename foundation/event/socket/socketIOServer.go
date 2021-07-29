package socket

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/foundation/system"
	"log"
	"net/http"
	"strconv"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
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
	server.On(socketio.OnConnection, onClientConnected)
	server.On(socketio.OnDisconnection, onClientDisconnected)
	// for status handling
	server.On("elastatus", func(socket *socketio.Channel) string {
		return s.GetStatus()
	})

	s.socket = server
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	s.state = data.Connected
	log.Println("Starting server @PORT ", constants.PORT)
	var TIMEOUT int64 = 10
	// step: try connecting
	go func() {
		elapsed := time.Now().Unix()
		for {
			err := http.ListenAndServe(":"+strconv.Itoa(constants.PORT), serveMux)
			if err == nil {
				break
			}
			// step: waiting for too long?
			diff := time.Now().Unix() - elapsed
			if diff > TIMEOUT {
				log.Fatal("Server error", err.Error())
				break
			}
			log.Println("Issue found, retrying...", err.Error())
			// sleep for a while
			time.Sleep(time.Millisecond * 500)
		}
	}()

	time.Sleep(time.Millisecond * 500)

	return nil
}

func (s *SocketIOServer) GetStatus() string {
	return system.GetStatus()
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
	println("SocketIOServer", "Broadcast", "room="+room, dataTransfer)
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
	clientCast := client.(*gosocketio.Channel)
	return clientCast.Ack(method, data, time.Second*constants.TIMEOUT)
}

func (s *SocketIOServer) SetStatus(status system.Status, data interface{}) error {
	system.SetStatus(string(status))
	return s.Broadcast(
		constants.SYSTEM_SERVICE_ID,
		constants.BCAST_SYSTEM_STATUS_CHANGED,
		statusData{status: string(status), data: data})
}

/// this closes the server
func (s *SocketIOServer) Close() {
	//go s.socket.()
	s.state = data.Disconnected
}
