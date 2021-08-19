package event

import (
	"ela/foundation/constants"
	"ela/foundation/system"
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
