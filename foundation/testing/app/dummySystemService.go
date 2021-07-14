package app

import (
	"ela/foundation/event/data"
)

type DummySystemService struct {
	pending data.ActionGroup
}

func NewDummy(pendingActions data.ActionGroup) *DummySystemService {
	return &DummySystemService{pending: pendingActions}
}

func (t *DummySystemService) RequestFor(action data.Action) (*data.Response, error) {
	return &data.Response{Value: t.pending.ToJson()}, nil
}

func (t *DummySystemService) Disconnect() error {
	return nil
}
