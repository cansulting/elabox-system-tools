package landing

import (
	"net/http"
	"os"

	econstants "github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/server"
)

const PORT = "80"

var serverhandler *server.Manager

const TIMEOUT = 10 // timeout for server initialization
var connected = 0

func Initialize(landingPagePath string) error {
	// step: serve event server
	serverhandler = &server.Manager{}
	serverhandler.Setup()

	// step: init web server
	fileserver := http.FileServer(http.Dir(landingPagePath))
	constants.Logger.Debug().Msg("Landing page path @" + landingPagePath)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		connected++
		url := r.URL.Path
		if _, err := os.Stat(landingPagePath + url); err == nil {
			fileserver.ServeHTTP(rw, r)
		} else {
			http.Redirect(rw, r, "/", http.StatusFound)
		}
	})
	if err := serverhandler.ListenAndServe(); err != nil {
		return err
	}
	serverhandler.EventServer.SetStatus(system.UPDATING, nil)
	return nil
}

func GetServer() *server.Manager {
	return serverhandler
}

// use to broadcast installation log to system
func BroadcastLog(msg string) {
	if serverhandler != nil && serverhandler.IsRunning() {
		serverhandler.EventServer.Broadcast(econstants.SYSTEM_SERVICE_ID, "log", msg)
	}
}

// use to shutdown the server
func Shutdown() error {
	constants.Logger.Info().Str("category", "networking").Msg("Shutting down event and server...")
	// close event server
	if serverhandler != nil {
		if err := serverhandler.Stop(); err != nil {
			return errors.SystemNew("Error closing connector.", err)
		}
	}
	return nil
}
