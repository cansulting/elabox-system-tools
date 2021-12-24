package dbutils

import (
	"database/sql"
	"errors"
	"reflect"
)

var objectDefs map[reflect.Type]*ObjectDef = make(map[reflect.Type]*ObjectDef)

// info and index of a field in an object
type FieldIndex struct {
	index  int
	sfield reflect.StructField
}

// set value to object
func (ins FieldIndex) setValue(targetObj reflect.Value, val interface{}) {
	newV := reflect.ValueOf(val).Elem()
	targetObj.Field(ins.index).Set(newV)
}

// get value of field from instance of an object
func (ins FieldIndex) getValue(targetObj reflect.Value) interface{} {
	insField := targetObj.Field(ins.index)
	return insField.Interface()
}

// this struct represents
type ObjectDef struct {
	fieldsI    map[string]FieldIndex // field definitions
	objType    reflect.Type          // type of target object
	scannedVal []interface{}         // use for temp
	dbColumns  []string              // columns on db
}

func GetObjectDef(t reflect.Type) *ObjectDef {
	def := objectDefs[t]
	if def != nil {
		return def
	}
	newDef := &ObjectDef{
		objType: t,
		fieldsI: make(map[string]FieldIndex),
	}
	objectDefs[t] = newDef
	return newDef
}

// convert db records to object instances
func (ins *ObjectDef) ToObj(rows *sql.Rows) ([]interface{}, error) {
	if ins.dbColumns == nil {
		columns, err := rows.Columns()
		ins.dbColumns = columns
		if err != nil {
			return nil, err
		}
	}
	// initialize tmp value
	if ins.scannedVal == nil {
		if err := ins.initTmpVals(rows); err != nil {
			return nil, err
		}
	}
	res := make([]interface{}, 0)
	// convert each row to an object
	for rows.Next() {
		err := rows.Scan(ins.scannedVal...)
		if err != nil {
			return nil, err
		}
		newInstance := reflect.New(ins.objType).Interface()
		res = append(res, newInstance)
		newObj := reflect.ValueOf(newInstance).Elem()
		// place scanned values to the object fields
		for i, v := range ins.scannedVal {
			field, ok := ins.getIndexByTag(ins.dbColumns[i])
			if !ok {
				continue
			}
			field.setValue(newObj, v)
		}
	}
	return res, nil
}

// initialize temp scannedVal to be use for scanning from db
func (ins *ObjectDef) initTmpVals(rows *sql.Rows) error {
	columnC := len(ins.dbColumns)
	ins.scannedVal = make([]interface{}, columnC)
	for i := range ins.scannedVal {
		field, ok := ins.getIndexByTag(ins.dbColumns[i])
		if !ok {
			ins.scannedVal[i] = new(interface{})
			continue
		}
		ins.scannedVal[i] = reflect.New(field.sfield.Type).Interface()
	}
	return nil
}

// get field info based on a tag
// return index value and true if returned successfully
func (ins *ObjectDef) getIndexByTag(tag string) (FieldIndex, bool) {
	val := ins.fieldsI[tag]
	if val.sfield.Index != nil {
		return val, true
	}

	// lookup via reflection
	fieldT := ins.objType.NumField()
	for i := 0; i < fieldT; i++ {
		f := ins.objType.Field(i)
		lookupTag, ok := f.Tag.Lookup("json")
		if ok && lookupTag == tag {
			f, _ := ins.objType.FieldByName(f.Name)
			newIns := FieldIndex{
				index:  i,
				sfield: f,
			}
			return newIns, true
		}
	}
	return FieldIndex{}, false
}

// use to get records based from object instance
// @param 1: the object to extract scannedVal
// @param 2: the resuable variable array that holds the output. doesnt need be empty
// @param 3: columns
func (ins *ObjectDef) ToRecords(obj interface{}, out *[]interface{}, columns ...string) error {
	if obj == nil {
		return errors.New("obj parameter shouldnt be nil")
	}
	if out == nil {
		return errors.New("out parameter that holds results should be initialized")
	}
	reflectedObj := reflect.ValueOf(obj).Elem()
	for i, column := range columns {
		tmpField, ok := ins.getIndexByTag(column)
		var value interface{} = nil
		if ok {
			value = tmpField.getValue(reflectedObj)
		}
		if len(*out)-1 < i {
			*out = append(*out, value)
		} else {
			(*out)[i] = value
		}
	}
	return nil
}
