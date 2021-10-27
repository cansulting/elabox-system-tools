// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package event

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"log"

	socketio "github.com/graarh/golang-socketio"
)

type statusData struct {
	status string
	data   interface{}
}

func (s *SocketIOServer) initStatus() {
	// for status handling
	s.socket.On("elastatus", func(socket *socketio.Channel) string {
		return s.GetStatus()
	})
}

// use to get status of system
func (s *SocketIOServer) GetStatus() string {
	return system.GetStatus()
}

// use to set the current status of system
func (s *SocketIOServer) SetStatus(status system.Status, data interface{}) error {
	log.Println("Server.SetStatus", status)
	system.SetStatus(string(status))
	if err := s.Broadcast(
		constants.SYSTEM_SERVICE_ID,
		constants.BCAST_SYSTEM_STATUS_CHANGED,
		statusData{status: string(status), data: data}); err != nil {
		log.Println("Server.SetStatus failure", err.Error())
		return err
	}
	return nil
}
