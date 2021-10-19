package app

import (
	_app "github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

func RunTestApp(controller *_app.Controller, pendingAction data.ActionGroup) {
	controller.RPC = NewDummy(pendingAction)
	_app.RunApp(controller)
}
