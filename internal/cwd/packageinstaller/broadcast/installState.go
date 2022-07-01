package broadcast

type InstallState string

const (
	INITIALIZING InstallState = "INIT"
	INPROGRESS   InstallState = "INPROGRESS"
	INSTALLED    InstallState = "INSTALLED"
	UNINSTALLED  InstallState = "UNINSTALLED"
)
