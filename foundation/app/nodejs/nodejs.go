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

func (n *Nodejs) Run() error {
	path := n.Config.GetNodejsDir() + "index.js"
	n.cmd = exec.Command("node", path)
	n.cmd.Dir = n.Config.GetNodejsDir()
	n.cmd.Stdout = n
	n.cmd.Stderr = n
	return n.cmd.Start()
}

func (n *Nodejs) Write(p []byte) (c int, err error) {
	log.Print(string(p))
	return len(p), nil
}

func (n *Nodejs) Stop() error {
	return n.cmd.Process.Kill()
}
