package protocol

import "ela/foundation/event/data"

type ActivityInterface interface {
	OnStart(action *data.Action) error
	IsRunning() bool
	OnEnd() error
}
