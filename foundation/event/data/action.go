package data

import "log"

type Action struct {
	// action id, which represents what action to make
	Id string `json:"id"`
	// optional. which specific package will handle this action.
	// if nothing was specified then look for any valid package that can carry out the action
	PackageId string `json:"pkid"`
	// optional. data which will be use to execute the action
	Data interface{} `json:"data"`
}

// convert Data to Action
func (a Action) DataToActionData() Action {
	val, ok := a.Data.(*Action)
	if ok {
		return *val
	} else {
		return Action{}
	}
}

// convert Action.data to int
func (a Action) DataToInt() int {
	switch a.Data.(type) {
	case int:
		return a.Data.(int)
	case float64:
		var f = a.Data.(float64)
		return int(f)
	case float32:
		var f = a.Data.(float32)
		return int(f)
	default:
		log.Panicln("Failed to concert Action to int ", a)
		return -1
	}
}

func (a Action) DataToString() string {
	return a.Data.(string)
}
