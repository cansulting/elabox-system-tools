package server

import (
	"context"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/server/config"
	"ela/server/event"
	"log"
	"net/http"
	"time"
)

type Manager struct {
	httpS         *http.Server
	EventServer   *event.SocketIOServer
	OnSystemEvent func(protocol.ClientInterface, data.Action) interface{} // callback for event server... specifically for system event usage
	running       bool
}

func (m *Manager) IsRunning() bool {
	return m.running
}

// setup the server
func (m *Manager) Setup() error {
	log.Println("Setting up web and event server...")
	m.httpS = &http.Server{Addr: ":" + config.PORT}
	// step: initialize event server
	eventS := &event.SocketIOServer{}
	m.EventServer = eventS
	if err := eventS.Open(); err != nil {
		return err
	}

	eventS.HandleSystemService(m.OnSystemEvent)
	return nil
}

// stop serving
func (m *Manager) Stop() error {
	m.running = false
	return m.httpS.Shutdown(context.TODO())
}

// start serving the server
func (m *Manager) ListenAndServe() {
	var TIMEOUT int64 = 10
	m.running = true
	// step: try connecting
	go func() {
		elapsed := time.Now().Unix()
		for m.running {
			err := m.httpS.ListenAndServe()
			if err == nil {
				break
			}
			// step: waiting for too long?
			diff := time.Now().Unix() - elapsed
			if diff > TIMEOUT {
				log.Println("Server manager error", err.Error())
				break
			}
			log.Println("Issue found, retrying...", err.Error())
			// sleep for a while
			time.Sleep(time.Millisecond * 500)
		}
		m.running = false
	}()
}
