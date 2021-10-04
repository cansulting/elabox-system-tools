package utils

import (
	"os"
	"os/exec"
	"path/filepath"

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

// restart the main system
func RestartSystem() error {
	pkc.Logger.Info().Msg("Restarting system...")
	// step: execute system binary
	systemPath := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
	cmd := exec.Command(systemPath)
	//cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Dir = filepath.Dir(systemPath)
	if err := cmd.Start(); err != nil {
		return errors.SystemNew("Restart system failed", err)
	}
	return cmd.Process.Release()
}

func TerminateSystem() error {
	pkc.Logger.Info().Msg("Terminating system...")
	// step: execute system binary
	systemPath := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
	if _, err := os.Stat(systemPath); err != nil {
		pkc.Logger.Warn().Msg("Terminate skipped. System is not installed.")
		return nil
	}
	cmd := exec.Command(systemPath, "terminate")
	cmd.Dir = filepath.Dir(systemPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
