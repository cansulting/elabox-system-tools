package data

type ConnectionType int

const (
	Disconnected = iota
	Connected
	Connecting
	Disconnecting
)
