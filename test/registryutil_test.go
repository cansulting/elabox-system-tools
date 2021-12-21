package main

import (
	"reflect"
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

func TestRowsToPackages(t *testing.T) {
	rows, err := util.SelectQuery("select * from packages")
	if err != nil {
		t.Error(err)
		return
	}
	objectD := util.GetObjectDef(reflect.TypeOf(data.PackageConfig{}))
	_, err = objectD.ToObj(rows)
	// _, err = util.RowsTo(rows, reflect.TypeOf(data.PackageConfig{}))
	if err != nil {
		return
	}
}
