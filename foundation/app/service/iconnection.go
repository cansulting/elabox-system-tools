package service

import (
	"ela/foundation/event/data"
)

type IConnection interface {
	RequestFor(action data.Action) (*data.Response, error)
	Disconnect() error
}
