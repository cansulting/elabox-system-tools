package main

import (
	"C"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
)

func main() {
	InitializePath()

	// install via commandline?
	args := os.Args
	if len(args) > 1 {
		startCommandline()
		return
	}
	// install via installer service
	var err error
	constants.AppController, err = app.NewController(&activity{}, nil)
	if err != nil {
		constants.Logger.Error().Err(err).Caller().Msg("Failed to initialize App Controller")
		panic(err)
	}
	if err := app.RunApp(constants.AppController); err != nil {
		constants.Logger.Error().Err(err).Caller().Msg("Failed running app")
		panic(err)
	}
}
