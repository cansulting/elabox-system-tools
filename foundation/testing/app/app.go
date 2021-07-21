package app

import (
	_app "ela/foundation/app"
	"ela/foundation/event/data"
)

func RunTestApp(controller *_app.Controller, pendingAction data.ActionGroup) {
	controller.RPC = NewDummy(pendingAction)
	_app.RunApp(controller)
}
