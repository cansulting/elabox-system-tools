package protocol

type ServiceInterface interface {
	OnStart() error
	IsRunning() bool
	OnEnd() error
}
