package data

type ServiceState int

const (
	ServiceAwake    = iota // service means currently running
	ServiceSleeping        // service is not executed
	ServiceStopped         // service was stopped
)
