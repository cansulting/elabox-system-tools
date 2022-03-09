package broadcast

import (
	econstants "github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/landing"
)

// use to broadcast installation log to system.
// directly call the socket server
func SystemLog(msg string) {
	serverhandler := landing.GetServer()
	if serverhandler != nil && serverhandler.IsRunning() {
		serverhandler.EventServer.Broadcast(econstants.SYSTEM_SERVICE_ID, "log", msg)
	}
}
