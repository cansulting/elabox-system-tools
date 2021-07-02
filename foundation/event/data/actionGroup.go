package data

type ActionGroup struct {
	Activities []Action
	Broadcasts []Action
}

func NewActionGroup() *ActionGroup {
	group := &ActionGroup{}
	group.Activities = make([]Action, 0, 3)
	group.Broadcasts = make([]Action, 0, 4)
	return group
}

func (app *ActionGroup) AddPendingActivity(action Action) {
	app.Activities = append(app.Activities, action)
}

func (app *ActionGroup) AddPendingBroadccast(action Action) {
	app.Broadcasts = append(app.Broadcasts, action)
}

func (app *ActionGroup) ClearAll() {
	app.Activities = app.Activities[:0]
	app.Broadcasts = app.Broadcasts[:0]
}
