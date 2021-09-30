package main

import (
	"ela/foundation/logger"
	pkconst "ela/internal/cwd/packageinstaller/constants"
	"ela/internal/cwd/packageinstaller/landing"
	"ela/internal/cwd/packageinstaller/pkg"
	"ela/internal/cwd/packageinstaller/utils"
	"os"
	"time"

	"github.com/rs/zerolog"
)

/*
	commandline.go
	Commandline version of installer
*/

type loggerHook struct {
}

func (i loggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	landing.BroadcastLog(msg)
}

func startCommandline() {
	println("Elabox Installer Commandline")
	println("type help or -h for arguments.")

	// step: commandline help?
	if IsArgExist("help") || IsArgExist("-h") {
		println("usage:")
		println("Example command <path to package> -r -s")
		println("-r - to restart the system")
		println("-s - this is system update")
		println("-l - create log file")
		println("-i - ignore custom installer")
		return
	}
	systemUpdate := IsArgExist("-s")
	// true if restarts system
	restartSystem := IsArgExist("-r") || systemUpdate
	args := os.Args
	targetPk := args[1]
	// step: load package
	content, err := pkg.LoadFromSource(targetPk)
	if err != nil {
		pkconst.Logger.Fatal().Err(err).Caller().Msg("Failed running commandline")
		return
	}
	// request for server broadcast
	if IsArgExist("-l") || systemUpdate {
		logger.SetHook(loggerHook{})
		pkconst.Logger = logger.GetInstance()
	}
	// step: we need clients to system update via ports
	if systemUpdate {
		startServer(content)
	} else {
		logger.ConsoleOut = false
	}
	// use custom installer or not?
	if IsArgExist("-i") || !content.HasCustomInstaller() {
		normalInstall(content)
	} else {
		if err := content.RunCustomInstaller(targetPk, true, "-i"); err != nil {
			pkconst.Logger.Fatal().Err(err).Caller().Msg("Failed running custom installer")
			return
		}
	}
	pkconst.Logger.Info().Msg("Installed success.")
	// step: stop listeners
	if systemUpdate {
		if err := landing.Shutdown(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Error shutting down.")
		}
	}
	// step: restart system
	if restartSystem {
		if err := utils.RestartSystem(); err != nil {
			pkconst.Logger.Fatal().Err(err)
			return
		}
		time.Sleep(time.Millisecond * 200)
		os.Exit(0)
	}
}

func normalInstall(content *pkg.Data) {
	// step: wait and make sure system was terminated. for system updates
	time.Sleep(time.Second)
	newInstall := NewInstaller(content, true)
	// step: start install
	if err := newInstall.Start(); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Failed reverting installer.")
		}
		pkconst.Logger.Fatal().Err(err).Stack()
		return
	}
	// step: post install
	if err := newInstall.Finalize(); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Failed reverting installer.")
		}
		pkconst.Logger.Fatal().Err(err).Stack()
		return
	}
}

// start installer server
func startServer(content *pkg.Data) {
	// retrieve landing page first
	landingDir, err := content.ExtractLandingPage()
	// is there a landing page?
	if err == nil {
		if err := landing.Initialize(landingDir); err != nil {
			pkconst.Logger.Fatal().Err(err).Stack().Msg("Unable to initialize server.")
			return
		}
	} else {
		pkconst.Logger.Error().Err(err).Caller().Msg("Failed extracting landing page. Skipping www listener")
	}
	// step: if theres a landing page. wait for user to connect to landing page before continuing
	//if landingDir != "" {
	//	landing.WaitForConnection()
	//}
}

func IsArgExist(arg string) bool {
	args := os.Args[1:]
	for _, _arg := range args {
		if arg == _arg {
			return true
		}
	}
	return false
}
