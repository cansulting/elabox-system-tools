package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	pkc "github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
)

// check if system is currently running
func isSystemRunning() bool {
	// step: connect to system
	if pkc.AppController == nil || pkc.AppController.RPC == nil {
		return false
	}
	return true
}

// start the main system
func StartSystem() error {
	pkc.Logger.Info().Msg("Starting up system...")
	// step: execute system binary
	systemPath := path.GetAppInstallLocation(constants.SYSTEM_SERVICE_ID, false) + "/" + constants.SYSTEM_SERVICE_ID
	cmd := exec.Command(systemPath)
	//cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Dir = filepath.Dir(systemPath)
	if err := cmd.Start(); err != nil {
		return errors.SystemNew("Restart system failed", err)
	}
	return cmd.Process.Release()
}

func TerminateSystem(delaySec int) error {
	// step: execute system binary
	systemPath := path.GetAppInstallLocation(constants.SYSTEM_SERVICE_ID, false) + "/" + constants.SYSTEM_SERVICE_ID
	if _, err := os.Stat(systemPath); err != nil {
		//pkc.Logger.Warn().Msg("Terminate skipped. System is not installed.")
		return nil
	}
	pkc.Logger.Info().Msg("terminating system in " + strconv.Itoa(delaySec) + " seconds")
	time.Sleep(time.Second * time.Duration(delaySec))
	cmd := exec.Command(systemPath, "terminate")
	cmd.Dir = filepath.Dir(systemPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// use to reboot the device
func Reboot() error {
	pkc.Logger.Info().Msg("Rebooting system...")
	cmd := exec.Command("reboot")
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Process.Release()
	return nil
}
