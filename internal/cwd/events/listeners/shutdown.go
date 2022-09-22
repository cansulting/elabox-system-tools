package listeners

import (
	"os"
	"syscall"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/ztrue/shutdown"
)

func ListenToShutdown() {
	shutdown.AddWithParam(func(sig os.Signal) {
		logger.GetInstance().Info().Msg("Shutdown signal received")
		system.SetEnv("ELASHUTDOWNSTATUS", "properly_shutdown")
	})
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}