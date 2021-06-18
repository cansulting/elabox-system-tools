package servicecenter

import (
	"log"
	"runtime"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type SocketIOConnectorClient struct {
	state  ConnectionType
	socket *gosocketio.Client
}

func (s *SocketIOConnectorClient) GetState() ConnectionType {
	return s.state
}

// implementation for connector client. let client send service requests
func (s *SocketIOConnectorClient) SendServiceRequest(serviceId string, data ActionData) (string, error) {
	log.Println("socketIOConnectorClient.SendServiceRequest", serviceId, data)
	return s.socket.Ack(serviceId, data, time.Second*TIMEOUT)
}

// implementation for connector client. let this client subscribe to specific room
func (s *SocketIOConnectorClient) Subscribe(
	subs SubscriptionData,
	callback interface{}) (string, error) {
	err := s.socket.On(subs.Action, callback)
	if err != nil {
		return "", err
	}
	return s.socket.Ack("subscribe", subs, time.Second*TIMEOUT)
}

func (s *SocketIOConnectorClient) Broadcast(action ActionData) error {
	_, err := s.socket.Ack("broadcast", action, time.Second*TIMEOUT)
	return err
}

func (s *SocketIOConnectorClient) Open() error {
	log.Println("Socket Connecting")
	if s.socket == nil {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	s.socket = nil

	var err error
	var c *gosocketio.Client
	for {
		client, err := gosocketio.Dial(
			gosocketio.GetUrl("localhost", PORT, false),
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 1)
		} else {
			c = client
			break
		}
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Println("Disconnected")
		if RECONNECT {
			log.Println("Reconnect")
			time.Sleep(time.Second * 1)
			s.Open()
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(" [x] initialize")
	s.socket = c
	return nil
}

func (s *SocketIOConnectorClient) Close() {
	s.socket.Close()
}
