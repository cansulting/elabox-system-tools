// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package data

import (
	"encoding/json"
	"log"
)

type ActionGroup struct {
	Activity *Action `json:"activity"`
	Service  *Action `json:"service"`
	//Broadcasts []Action
}

func NewActionGroup() *ActionGroup {
	group := &ActionGroup{}
	//group.Activities = make([]Action, 0, 3)
	//group.Broadcasts = make([]Action, 0, 4)
	return group
}
func NewActionGroupFromMap(data map[string]interface{}) *ActionGroup {
	group := &ActionGroup{}
	//group.Activities = make([]Action, 0, 3)
	//group.Broadcasts = make([]Action, 0, 4)
	if val, ok := data["activity"]; ok {
		activity := val.(map[string]interface{})
		actionId := ""
		packageId := ""
		var acData interface{}
		if val, ok := activity["id"]; ok {
			actionId = val.(string)
		}
		if val, ok := activity["packageId"]; ok {
			packageId = val.(string)
		}
		if val, ok := activity["data"]; ok {
			acData = val
		}
		tmp := NewAction(actionId, packageId, acData)
		group.Activity = &tmp
	}
	return group
}

func (app *ActionGroup) AddPendingActivity(action *Action) {
	app.Activity = action
	//app.Activities = append(app.Activities, action)
}

func (app *ActionGroup) AddPendingService(action *Action) {
	app.Service = action
}

func (app *ActionGroup) AddPendingBroadccast(action Action) {
	//app.Broadcasts = append(app.Broadcasts, action)
}

func (app *ActionGroup) ClearAll() {
	//app.Activities = app.Activities[:0]
	//app.Broadcasts = app.Broadcasts[:0]
}

func (app *ActionGroup) ToJson() string {
	res, err := json.Marshal(app)
	if err != nil {
		log.Println("ActionGroup.ToJson() failed to marshal")
		return ""
	}
	return string(res)
}

/*
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
*/
