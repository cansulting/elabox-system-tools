package appman

import (
	"os/exec"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
)

/*
	Structure that handles Nodejs runtime
*/
type Nodejs struct {
	Config  *data.PackageConfig
	cmd     *exec.Cmd
	running bool
}

// return true if currently running
func (n *Nodejs) IsRunning() bool {
	return n.running
}

// start running node js
func (n *Nodejs) Run() error {
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
		return errors.SystemNew("Failed running nodejs "+n.Config.PackageId, err)
	}
	return nil
	//n.Stop()
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
	return n.cmd.Process.Kill()
}
