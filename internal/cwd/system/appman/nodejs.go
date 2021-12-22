package appman

import (
	"os/exec"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
)

/*
	Structure that handles Nodejs runtime
*/
type Nodejs struct {
	Config           *data.PackageConfig
	RestartOnFailure bool // true if nodejs will automatically restart upon failure
	cmd              *exec.Cmd
	running          bool
}

// return true if currently running
func (n *Nodejs) IsRunning() bool {
	return n.running
}

// start running node js
func (n *Nodejs) Run() error {
	if n.running {
		return nil
	}
	n.running = true
	defer func() { n.running = false }()

	path := n.Config.GetNodejsDir() + "/index.js"
	cmd := exec.Command("node", path)
	n.cmd = cmd
	cmd.Dir = n.Config.GetNodejsDir()
	cmd.Stdout = n
	cmd.Stderr = n
	if err := cmd.Start(); err != nil {
		return errors.SystemNew("Failed Starting nodejs "+n.Config.PackageId, err)
	}
	if err := cmd.Wait(); err != nil {
		if n.running {
			return n.onRunFailure(err)
		}
	}
	return nil
	//n.Stop()
}

// callback when nodejs failed unexpectedly while running.
func (n *Nodejs) onRunFailure(err error) error {
	global.Logger.Error().Err(err).Msg("Unexpected termination for " + n.Config.PackageId)
	if n.RestartOnFailure {
		return n.Restart()
	} else {
		return errors.SystemNew("Failed running nodejs "+n.Config.PackageId, err)
	}
}

// use to restart the current node js
func (n *Nodejs) Restart() error {
	if err := n.Stop(); err != nil {
		global.Logger.Error().Err(err).Msg("Failed stopping nodejs " + n.Config.PackageId)
	}
	global.Logger.Log().Msg("Restarting nodejs " + n.Config.PackageId)
	return n.Run()
}

// callback when system has log
func (n *Nodejs) Write(data []byte) (int, error) {
	print(string(data))
	return len(data), nil
}

func (n *Nodejs) Stop() error {
	if !n.running {
		return nil
	}
	n.running = false
	return n.cmd.Process.Kill()
}
