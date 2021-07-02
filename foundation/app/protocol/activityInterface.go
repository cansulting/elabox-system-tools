package protocol

type ActivityInterface interface {
	OnStart() error
	IsRunning() bool
	OnEnd() error
}
