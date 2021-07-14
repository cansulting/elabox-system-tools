package socket

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
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
	action string,
	callback interface{}) error {
	err := s.socket.On(action, callback)
	if err != nil {
		return err
	}
	return nil
}

func (s *SocketIOClient) Broadcast(event string, data interface{}) error {
	_, err := s.socket.Ack(event, data, time.Second*constants.TIMEOUT)
	return err
}

// use to connect to local app server
// @timeout: time in seconds it will timeout. @timeout > 0 to apply timeout
func (s *SocketIOClient) Open(timeout int16) error {
	log.Println("Socket Connecting")
	if s.socket == nil {
		runtime.GOMAXPROCS(1 /*runtime.NumCPU()*/)
	}

	s.socket = nil

	var err error

	// step: try to establish connection
	elapsedTimeout := int16(0)
	var c *gosocketio.Client
	for {
		client, err := gosocketio.Dial(
			gosocketio.GetUrl("localhost", constants.PORT, false),
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			if timeout > 0 {
				if elapsedTimeout >= timeout {
					return errors.SystemNew("Timeout", nil)
				}
				elapsedTimeout++
			}
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
			s.Open(-1)
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
