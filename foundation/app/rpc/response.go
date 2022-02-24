// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// represents data from remote request

package rpc

import (
	"encoding/json"
	"errors"

	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

type ResponseMessage struct {
	Code    float32 `json:"code"`
	Message string  `json:"message"`
}

type Response struct {
	Value interface{}
}

func (inst Response) toStringDecoded() (string, error) {
	return DecodeResponse(inst.Value.(string))
}

func (inst Response) ToSimpleResponse() (ResponseMessage, error) {
	val := &ResponseMessage{}
	if err := inst.ParseJson(val); err != nil {
		return ResponseMessage{}, err
	}
	return *val, nil
}

func (r *Response) ParseJson(obj interface{}) error {
	if r.Value != nil {
		strVal := []byte(r.ToString())
		if err := json.Unmarshal(strVal, obj); err != nil {
			return err
		}
		return nil
	}
	return errors.New("cannot parse empty value")
}

func (r *Response) ToActionGroup() *data.ActionGroup {
	actiong := data.NewActionGroup()
	r.ParseJson(&actiong)
	return actiong
}

func (r *Response) ToString() string {
	if dec, err := r.toStringDecoded(); err == nil {
		return dec
	}
	return ""
}
