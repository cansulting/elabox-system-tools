package utils

import (
	"os/exec"
	"strconv"
)

// use to allow specific port via ufw command
func AllowPort(port int) error {
	cmd := exec.Command("ufw", "allow", strconv.Itoa(port))
	return cmd.Start()
}
