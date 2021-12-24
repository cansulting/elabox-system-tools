package main

import (
	"reflect"
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/registry/dbutils"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

func TestDBSerialization(t *testing.T) {
	rows, err := util.SelectQuery("select * from packages")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("testing serializing objects from records")
	columns, _ := rows.Columns()
	objectD := dbutils.GetObjectDef(reflect.TypeOf(data.PackageConfig{}))
	objs, err := objectD.ToObj(rows)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("testing deserialize object to record")
	values := make([]interface{}, 1)
	if err := objectD.ToRecords(objs[0], &values, columns...); err != nil {
		t.Error("failed creating records based from object", err)
		return
	}
}

func TestDBContruct() {

}
