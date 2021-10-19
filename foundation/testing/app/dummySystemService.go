package app

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/service"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

type DummyRPC struct {
	pending data.ActionGroup
}

func NewDummy(pendingActions data.ActionGroup) *DummyRPC {
	return &DummyRPC{pending: pendingActions}
}

func (t *DummyRPC) Call(packageId string, action data.Action) (*data.Response, error) {
	return &data.Response{Value: t.pending.ToJson()}, nil
}

func (t *DummyRPC) CallSystem(action data.Action) (*data.Response, error) {
	return &data.Response{Value: t.pending.ToJson()}, nil
}

func (t *DummyRPC) OnRecieved(event string, onServiceResponse service.ServiceDelegate) {

}

func (t *DummyRPC) Close() error {
	return nil
}
