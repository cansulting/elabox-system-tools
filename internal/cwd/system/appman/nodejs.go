package appman

import (
	"ela/foundation/app/data"
	"log"
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

func (n *Nodejs) Run() {
	n.running = true
	path := n.Config.GetNodejsDir() + "/index.js"
	cmd := exec.Command("node", path)
	n.cmd = cmd
	cmd.Dir = n.Config.GetNodejsDir()
	cmd.Stdout = n
	cmd.Stderr = n
	if err := cmd.Start(); err != nil {
		log.Println("System.Nodejs:", n.Config.PackageId, "ERROR", err)
		return
	}
	if err := cmd.Wait(); err != nil {
		log.Println("System.Nodejs:", n.Config.PackageId, "ERROR", err)
	}
	n.running = false
	//n.Stop()
}

// callback when system has log
func (n *Nodejs) Write(data []byte) (int, error) {
	log.Print(string(data))
	return len(data), nil
}

func (n *Nodejs) Stop() error {
	if !n.running {
		return nil
	}
	return n.cmd.Process.Kill()
}
