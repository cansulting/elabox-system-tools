package appman

import (
	"ela/foundation/app/data"
	"ela/foundation/errors"
	"os/exec"
)

/*
	Structure that handles Nodejs runtime
*/
type Nodejs struct {
	Config  *data.PackageConfig
	cmd     *exec.Cmd
	running bool
}

// start running node js 
func (n *Nodejs) Run() error {
	n.running = true
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
	n.running = false
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
