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
