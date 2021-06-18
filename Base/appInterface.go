package base

type AppInterface interface {
	IsRunning() bool
	GetService() ServiceInterface
	OnStart() error
	OnEnd() error
}
