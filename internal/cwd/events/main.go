package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/events/listeners"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
)

func main() {
	logger.Init("events")
	logger.GetInstance().Info().Msg("Starting events listener")
	listeners.ListenToShutdown()
	var err error
	constants.AppController, err = app.NewController(&activity{}, nil)
	if err != nil {
		constants.Logger.Error().Err(err).Caller().Msg("Failed to initialize App Controller")
		panic(err)
	}
	if err := app.RunApp(constants.AppController); err != nil {
		constants.Logger.Error().Err(err).Stack().Msg("Failed running app")
		panic(err)
	}	
}