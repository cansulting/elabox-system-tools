package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/server/config"
	"github.com/cansulting/elabox-system-tools/server/event"
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

// start server requests handling
func (m *Manager) ListenAndServe() error {
	var TIMEOUT int64 = 20
	m.running = true

	addr := m.httpS.Addr
	if addr == "" {
		addr = ":http"
	}
	var ln net.Listener
	var err error
	elapsed := time.Now().Unix()

	// try to start listening. until it become available
	for {
		ln, err = net.Listen("tcp", addr)
		if err == nil {
			break
		}
		// step: waiting for too long?
		diff := time.Now().Unix() - elapsed
		if diff > TIMEOUT {
			return errors.New("Server initialization timeout. Stop existing process that uses port before continuing. Inner error - " + err.Error())
		}
		logger.GetInstance().Error().Err(err).Str("category", "networking").Caller().Msg("Issue found, retrying...")
		// sleep for a while
		time.Sleep(time.Second)
	}
	// step: try connecting
	if err == nil {
		go func() {
			err := m.httpS.Serve(ln)
			if err != nil && m.running {
				logger.GetInstance().Error().Err(err).Str("category", "networking").Caller().Msg("Server serve failed..")
			}
			m.running = false
		}()
	}
	return nil
}
