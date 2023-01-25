package main

import (
	"os"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/broadcast"
	pkconst "github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/landing"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/pkg"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/sysinstall.go"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/sysupgrade"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/utils"

	"github.com/rs/zerolog"
)

/*
	commandline.go
	Commandline version of installer
*/

type loggerHook struct {
}

func (i loggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	broadcast.SystemLog(msg)
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
		println("-u -uninstall package")
		return
	}
	pk := os.Args[1]
	if IsArgExist("-u") {
		processUninstallCommand(pk, IsArgExist("-l"))
	} else {
		processInstallCommand(pk, IsArgExist("-r"), IsArgExist("-l"), !IsArgExist("-i"))
	}

}

func processInstallCommand(targetPk string, restart bool, logging bool, runCustomInstaller bool) {
	// step: load package
	content, err := pkg.LoadFromSource(targetPk)
	if err != nil {
		pkconst.Logger.Fatal().Err(err).Caller().Msg("Failed running commandline")
		return
	}
	systemUpdate := false
	if content.Config.PackageId == constants.SYSTEM_SERVICE_ID {
		systemUpdate = true
	}
	// step: request for server broadcast
	if logging || systemUpdate {
		logger.SetHook(loggerHook{})
		pkconst.Logger = logger.GetInstance()
	}
	// true if this is parent system update. parent system update means this will execute custom installer if there is
	parentSystemUpdate := false
	if systemUpdate && runCustomInstaller {
		parentSystemUpdate = true
	}

	// step: we need clients to system update
	if parentSystemUpdate {
		// step: terminate system
		// if err := utils.TerminateSystem(pkconst.TERMINATE_TIMEOUT); err != nil {
		// 	pkconst.Logger.Debug().Err(err).Caller().Msg("failed terminating system")
		// }
		startServer(content)
		// step: check if theres a last failed installation
		lastState := sysinstall.GetLastState()
		if lastState == sysinstall.SUCCESS {
			if err := sysinstall.MarkInprogress(); err != nil {
				pkconst.Logger.Error().Err(err).Stack().Caller().Msg("Failed to mark installation as inprogress. install aborted.")
				return
			}
		} else {
			pkconst.Logger.Info().Msg("Last installation was unsuccessfull, resuming ...")
		}
	} else {
		logger.ConsoleOut = false
		if systemUpdate {
			sysupgrade.CheckAndUpgrade(int(content.Config.Build))
		}
	}
	// step: use custom installer or not?
	if !runCustomInstaller || !content.HasCustomInstaller() {
		normalInstall(content)
	} else {
		if err := content.RunCustomInstaller(targetPk, true, "-i"); err != nil {
			pkconst.Logger.Fatal().Err(err).Caller().Msg("Failed running custom installer")
			return
		}
	}
	// step: and stop listeners
	if parentSystemUpdate {
		// mark as success
		if err := sysinstall.MarkSuccess(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Failed installation.")
			return
		}
		// shutdown
		if err := landing.Shutdown(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Error shutting down.")
		}
	}
	pkconst.Logger.Info().Msg("Installed success.")
	// step: restart system
	if parentSystemUpdate {
		if restart {
			if err := utils.Reboot(5); err != nil {
				pkconst.Logger.Fatal().Err(err)
				return
			}
		} else {
			if err := utils.StartSystem(); err != nil {
				pkconst.Logger.Fatal().Err(err)
				return
			}
		}
		time.Sleep(time.Millisecond * 200)
		os.Exit(0)
	}
}

func normalInstall(content *pkg.Data) {
	// step: wait and make sure system was terminated. for system updates
	time.Sleep(time.Second)
	newInstall := NewInstaller(content, true, false)
	// step: start install
	if err := newInstall.Start(); err != nil {
		pkconst.Logger.Error().Err(err).Stack().Msg("Failed installation. Reason: " + err.Error())
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Failed reverting installer.")
		}
		panic("Failed installing @normalInstall()")
	}
	// step: post install
	if err := newInstall.Finalize(); err != nil {
		pkconst.Logger.Error().Err(err).Caller().Stack().Msg("Failed installation. Reason: " + err.Error())
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("Failed reverting installer.")
		}
		panic("Failed installing @normalInstall()")
	}
}

// package uninstall

func processUninstallCommand(targetPk string, logging bool) error {
	if logging {
		logger.SetHook(loggerHook{})
		pkconst.Logger = logger.GetInstance()
	}
	if err := utils.UninstallPackage(targetPk, false, false, false); err != nil {
		pkconst.Logger.Error().Err(err).Caller().Msg("unable to uninstall package " + targetPk)
		return err
	}
	return nil
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
