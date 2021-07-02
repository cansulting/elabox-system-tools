package socket

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"log"
	"runtime"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type SocketIOClient struct {
	state  data.ConnectionType
	socket *gosocketio.Client
}

func (s *SocketIOClient) GetState() data.ConnectionType {
	return s.state
}

// implementation for connector client. let client send service requests
func (s *SocketIOClient) SendServiceRequest(serviceId string, action data.Action) (string, error) {
	log.Println("socketIOConnectorClient.SendServiceRequest", serviceId, action)
	return s.socket.Ack(serviceId, action, time.Second*constants.TIMEOUT)
}

// implementation for connector client. let this client subscribe to specific room
func (s *SocketIOClient) Subscribe(
	subs data.Subscription,
	callback interface{}) (string, error) {
	err := s.socket.On(subs.Action, callback)
	if err != nil {
		return "", err
	}
	return s.socket.Ack("subscribe", subs, time.Second*constants.TIMEOUT)
}

func (s *SocketIOClient) Broadcast(action data.Action) error {
	_, err := s.socket.Ack("broadcast", action, time.Second*constants.TIMEOUT)
	return err
}

// use to connect to local app server
func (s *SocketIOClient) Open() error {
	log.Println("Socket Connecting")
	if s.socket == nil {
		runtime.GOMAXPROCS(1 /*runtime.NumCPU()*/)
	}

	s.socket = nil

	var err error

	// step: try to establish connection
	var c *gosocketio.Client
	for {
		client, err := gosocketio.Dial(
			gosocketio.GetUrl("localhost", constants.PORT, false),
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 1)
		} else {
			c = client
			break
		}
	}

	// step: initialize disconnection event
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Println("Disconnected")
		if constants.RECONNECT {
			log.Println("Reconnect")
			time.Sleep(time.Second * 1)
			s.Open()
		}
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	// step: initialize connection event
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(" [x] initialize")
	s.socket = c
	return nil
}

func (s *SocketIOClient) Close() {
	s.socket.Close()
}
