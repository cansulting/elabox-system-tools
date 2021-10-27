// this file provides util for handling flag for system installation state
// this provides info if system was successfully installed or not

package sysinstall

import "os"

type eInstallState int

const (
	SUCCESS eInstallState = iota
	FAILED
)

// get the last state of installation
// return the last state and the old package info
func GetLastState() eInstallState {
	state := SUCCESS
	if HasOldPackage() {
		state = FAILED
	}
	return state
}

// system installation will be mark as inprogress
func MarkInprogress() error {
	return CreateOldPackageInfo()
}

func MarkSuccess() error {
	if _, err := os.Stat(OLD_PK); err != nil {
		return nil
	}
	return os.Remove(OLD_PK)
}
