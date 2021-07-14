package utils

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/internal/cwd/packageinstaller/global"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// check if system is currently running
func isSystemRunning() bool {
	// step: connect to system
	if global.AppController == nil || global.AppController.RPC == nil {
		return false
	}
	return true
}

// use restart the main system
func RestartSystem() error {
	// step: skip if current system is already running
	if !isSystemRunning() {
		return nil
	}
	log.Println("Restarting system...")
	// step: execute system binary
	systemPath := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
	cmd := exec.Command(systemPath)
	cmd.Dir = filepath.Dir(systemPath)
	if err := cmd.Start(); err != nil {
		return errors.SystemNew("Restart system failed", err)
	}
	time.Sleep(time.Second * 3)
	os.Exit(1)
	return nil
}
