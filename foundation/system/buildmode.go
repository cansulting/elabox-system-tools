package system

type BuildModeType string

const (
	DEBUG   BuildModeType = "DEBUG"
	STAGING BuildModeType = "STAGING"
	RELEASE BuildModeType = "RELEASE"
)
