package nodejs

import (
	"ela/foundation/app/data"
	"log"
	"os/exec"
)

type Nodejs struct {
	Config *data.PackageConfig
	cmd    *exec.Cmd
}

func (n *Nodejs) Run() {
	path := n.Config.GetNodejsDir() + "/index.js"
	n.cmd = exec.Command("node", path)
	n.cmd.Dir = n.Config.GetNodejsDir()
	output, err := n.cmd.CombinedOutput()
	log.Println("System.Nodejs:", n.Config.PackageId, string(output))
	if err != nil {
		log.Println("System.Nodejs:", n.Config.PackageId, "ERROR", err)
	}
	n.Stop()
}

func (n *Nodejs) Stop() error {
	return n.cmd.Process.Kill()
}
