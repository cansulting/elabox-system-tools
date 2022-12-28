package broadcast

type PkInstallerState string

// package installer broacast actions
const (
	INIT        PkInstallerState = "INIT"
	INPROGRESS  PkInstallerState = "INPROGRESS"
	INSTALLED   PkInstallerState = "INSTALLED"
	UNINSTALLED PkInstallerState = "UNINSTALLED"
)

// callback when theres progress update from installer
var OnInstallerProgress = func(string, int) {}

// callback when the package state changed
// @param 1: package id
// @param 2: the installer state
var OnInstallerStateChanged = func(string, PkInstallerState) {}

// callback when theres an error with package
// @param 1: packageid
// @param 2: error code
// @param 3: error msg
var OnInstallerError = func(string, int, string) {}
