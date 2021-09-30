package server

import (
	"context"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/foundation/logger"
	"ela/server/config"
	"ela/server/event"
	"net/http"
	"runtime"
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
	logger.GetInstance().Debug().Str("category", "networking").Msg("Setting up web and event server...")
	runtime.GOMAXPROCS(runtime.NumCPU())
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
				logger.GetInstance().Error().Err(err).Str("category", "networking").Caller().Msg("Server manager error.")
				break
			}
			logger.GetInstance().Error().Err(err).Str("category", "networking").Caller().Msg("Issue found, retrying...")
			// sleep for a while
			time.Sleep(time.Millisecond * 500)
		}
		m.running = false
	}()
}
