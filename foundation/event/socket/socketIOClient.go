// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This class implements connectorClient interface with the use of socketio

package socket

import (
	"log"
	"runtime"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"

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
func (s *SocketIOClient) SendSystemRequest(event string, action data.Action) (string, error) {
	//println("socketIOConnectorClient.SendSystemRequest", event)
	return s.socket.Ack(event, action, time.Second*constants.TIMEOUT)
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
