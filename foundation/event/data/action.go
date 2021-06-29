package data

import "log"

type Action struct {
	Action string `json:"action"`
	// optional value
	AppId string      `json:"appid"`
	Data  interface{} `json:"data"`
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
