package protocol

type AppInterface interface {
	IsRunning() bool
	GetService() ServiceInterface
	OnStart() error
	OnEnd() error
}
