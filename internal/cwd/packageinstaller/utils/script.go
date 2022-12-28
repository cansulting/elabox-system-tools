package utils

import (
	"io"
	"os"
	"os/exec"
)

const sh = "/bin/bash"

func ExecScript(script string, dir string, writer io.Writer) error {
	cmd := exec.Command(sh, script)
	if _, err := os.Stat(dir); err == nil {
		cmd.Dir = dir
	}
	if writer != nil {
		cmd.Stdout = writer
		cmd.Stderr = writer
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
