package servicecenter

import "log"

type ActionData struct {
	Action string `json:"action"`
	// optional value
	AppId string      `json:"appid"`
	Data  interface{} `json:"data"`
}

// convert Data to ActionData
func (a ActionData) DataToActionData() ActionData {
	val, ok := a.Data.(*ActionData)
	if ok {
		return *val
	} else {
		return ActionData{}
	}
}

// convert ActionData.data to int
func (a ActionData) DataToInt() int {
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
		log.Panicln("Failed to concert ActionData to int ", a)
		return -1
	}
}

type SubscriptionData struct {
	Action string `json:"action"`
	AppId  string `json:"appId"`
}
