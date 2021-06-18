package base

type ServiceInterface interface {
	OnStart() error
	IsRunning() bool
	OnEnd() error
}
