package system

import "os"

type Status string

const (
	STOPPED     = "inactive"
	RUNNING     = "active"
	BOOTING     = "booting"
	INIT_UPDATE = "init_update"
	UPDATING    = "updating"
)

func GetStatus() string {
	return os.Getenv("elastatus")
}

func SetStatus(status string) {
	os.Setenv("elastatus", status)
}
