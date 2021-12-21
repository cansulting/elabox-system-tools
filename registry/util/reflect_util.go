package util

import (
	"database/sql"
	"reflect"
)

var objectDefs map[reflect.Type]*ObjectDef = make(map[reflect.Type]*ObjectDef)

type FieldIndex struct {
	index  int
	sfield reflect.StructField
}

type ObjectDef struct {
	fieldsI    map[string]FieldIndex // field definitions
	otype      reflect.Type          // type of target object
	values     []interface{}         // use for temp
	rowColumns []string              // columns on db
}

func GetObjectDef(t reflect.Type) *ObjectDef {
	def := objectDefs[t]
	if def != nil {
		return def
	}
	newDef := &ObjectDef{
		otype:   t,
		fieldsI: make(map[string]FieldIndex),
	}
	objectDefs[t] = newDef
	return newDef
}

func (ins *ObjectDef) ToObj(rows *sql.Rows) ([]interface{}, error) {
	// initialize tmp value
	if ins.values == nil {
		if err := ins.initTmpVals(rows); err != nil {
			return nil, err
		}
	}
	if ins.rowColumns == nil {
		columns, err := rows.Columns()
		ins.rowColumns = columns
		if err != nil {
			return nil, err
		}
	}
	res := make([]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(ins.values...)
		if err != nil {
			return nil, err
		}
		newInstance := reflect.New(ins.otype).Interface()
		res = append(res, newInstance)
		reflectVal := reflect.ValueOf(newInstance).Elem()
		for i, v := range ins.values {

			field, ok := ins.getIndexByTag(ins.rowColumns[i])
			if !ok {
				continue
			}
			val := reflect.ValueOf(v).Elem()
			reflectVal.Field(field.index).Set(val)
		}
	}
	return res, nil
}

func (ins *ObjectDef) initTmpVals(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnC := len(columns)
	ins.values = make([]interface{}, columnC)
	for i := range ins.values {
		field, ok := ins.getIndexByTag(columns[i])
		if !ok {
			ins.values[i] = new(interface{})
			continue
		}
		ins.values[i] = reflect.New(field.sfield.Type).Interface()
	}
	return nil
}

func (ins *ObjectDef) getIndexByTag(tag string) (FieldIndex, bool) {
	val := ins.fieldsI[tag]
	if val.sfield.Index != nil {
		return val, true
	}

	fieldT := ins.otype.NumField()
	for i := 0; i < fieldT; i++ {
		f := ins.otype.Field(i)
		lookupTag, ok := f.Tag.Lookup("json")
		if ok && lookupTag == tag {
			f, _ := ins.otype.FieldByName(f.Name)
			newIns := FieldIndex{
				index:  i,
				sfield: f,
			}
			return newIns, true
		}
	}
	return FieldIndex{}, false
}
