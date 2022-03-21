package broadcast

type InstallState string

const (
	INITIALIZING InstallState = "INIT"
	INPROGRESS   InstallState = "INPROGRESS"
	SUCCESS      InstallState = "SUCCESS"
)
