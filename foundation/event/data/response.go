// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// represents data from remote request

package data

import (
	"encoding/json"
	"errors"
)

type Response struct {
	Value interface{}
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

func (r *Response) ToActionGroup() *ActionGroup {
	actiong := NewActionGroup()
	r.ParseJson(&actiong)
	return actiong
}

func (r *Response) ToString() string {
	return r.Value.(string)
}
