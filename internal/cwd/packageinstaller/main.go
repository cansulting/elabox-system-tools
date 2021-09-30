package main

import (
	"C"
	"ela/foundation/app"
	"ela/internal/cwd/packageinstaller/constants"
	"os"
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
		constants.Logger.Fatal().Err(err).Caller().Msg("Failed to initialize App Controller")
		return
	}
	if err := app.RunApp(constants.AppController); err != nil {
		constants.Logger.Fatal().Err(err).Caller().Msg("Failed running app")
	}
}
