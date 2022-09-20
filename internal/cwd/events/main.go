package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/events/listeners"
)

func main() {
	logger.Init("events")
	logger.GetInstance().Info().Msg("Starting events listener")
	listeners.ListenToShutdown()
}