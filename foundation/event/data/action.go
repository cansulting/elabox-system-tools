package data

import (
	"ela/foundation/errors"
	"encoding/json"
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Action struct {
	// action id, which represents what action to make
	Id string `json:"id"`
	// optional. which specific package will handle this action.
	// if nothing was specified then look for any valid package that can carry out the action
	PackageId string `json:"packageId"`
	// optional. data which will be use to execute the action
	Value interface{} `json:"data"`
	//valueAction *Action     `json:"-"`
}

func NewAction(id string, packageId string, data interface{}) Action {
	action := Action{
		Id:        id,
		PackageId: packageId,
	}
	action.Value = convertData(data)
	return action
}

func NewActionById(id string) Action {
	return NewAction(id, "", nil)
}

// convert Data to Action
func (a *Action) DataToActionData() (Action, error) {
	//if a.valueAction != nil {
	//	return *a.valueAction
	//}
	action := Action{}
	if a.Value == nil {
		return action, nil
	}
	switch a.Value.(type) {
	case string:
		strObj := a.DataToString()
		if err := json.Unmarshal([]byte(strObj), &action); err != nil {
			return action, errors.SystemNew("Action.valueToActionData failed to convert to Action", err)
		}
		break
	case map[string]interface{}:
		mapstructure.Decode(a.Value, &action)
		break
	}

	//a.valueAction = &action
	return action, nil
}

// convert Action.Value to int
func (a *Action) DataToInt() int {
	switch a.Value.(type) {
	case int:
		return a.Value.(int)
	case float64:
		var f = a.Value.(float64)
		return int(f)
	case float32:
		var f = a.Value.(float32)
		return int(f)
	default:
		log.Panicln("Failed to concert Action to int ", reflect.TypeOf(a.Value))
		return -1
	}
}

func (a *Action) DataToString() string {
	if a.Value != nil {
		return a.Value.(string)
	}
	return ""
}

func convertData(data interface{}) interface{} {
	if data != nil {
		switch data.(type) {
		case Action:
			tmpd := data.(Action)
			return tmpd.ToJson()
			//case ActionGroup:
			//	a.Value = data.(*ActionGroup).ToJson()
		default:
			return data
		}
	}
	return nil
}

func (a *Action) ToJson() string {
	res, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(res)
}
