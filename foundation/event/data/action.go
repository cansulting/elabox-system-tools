package data

// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// this file provides funnction for generating response data for RPC
// response data is json string which will be encoded to Base64 before transmitting

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/cansulting/elabox-system-tools/foundation/errors"

	"github.com/mitchellh/mapstructure"
)

type Action struct {
	// action id, which represents what action to make
	Id string `json:"id"`
	// optional. which specific package will handle this action.
	// if nothing was specified then look for any valid package that can carry out the action
	PackageId string `json:"packageId"`
	// optional. data which will be use to execute the action
	Data interface{} `json:"data"`
	//valueAction *Action     `json:"-"`
}

func NewAction(id string, packageId string, data interface{}) Action {
	action := Action{
		Id:        id,
		PackageId: packageId,
	}
	action.Data = convertData(data)
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
	if a.Data == nil {
		return action, nil
	}
	switch a.Data.(type) {
	case string:
		strObj := a.DataToString()
		if err := json.Unmarshal([]byte(strObj), &action); err != nil {
			return action, errors.SystemNew("Action.valueToActionData failed to convert to Action", err)
		}
		break
	case map[string]interface{}:
		mapstructure.Decode(a.Data, &action)
		break
	}

	//a.valueAction = &action
	return action, nil
}

// convert Action.Data to int
func (a *Action) DataToInt() int {
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
		log.Panicln("Failed to concert Action to int ", reflect.TypeOf(a.Data))
		return -1
	}
}

func (a *Action) DataToString() string {
	if a.Data != nil {
		return a.Data.(string)
	}
	return ""
}


func (a *Action) DataToMap(dst map[string]interface{}) (map[string]interface{}, error) {

	switch a.Data.(type) {
	case map[string]interface{}:
		return a.Data.(map[string]interface{}), nil
	}
	str := a.DataToString()
	if str != "" {
		if dst == nil {
			dst = make(map[string]interface{})
		}
		if err := json.Unmarshal([]byte(str), &dst); err != nil {
			return dst, err
		}
	}
	return nil, nil
}

func convertData(data interface{}) interface{} {
	if data != nil {
		switch data.(type) {
		case Action:
			tmpd := data.(Action)
			return tmpd.ToJson()
			//case ActionGroup:
			//	a.Data = data.(*ActionGroup).ToJson()
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
