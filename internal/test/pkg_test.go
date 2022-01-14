package main

import (
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
)

// test package loading
func TestLoadPackage(t *testing.T) {
	pkg := data.DefaultPackage()
	if err := pkg.LoadFromSrc("testpkg.json"); err != nil {
		t.Error(err)
		return
	}
	property, issue := pkg.GetIssue()
	if property != "" {
		t.Error("Theres an issue with", property, issue)
		return
	}

	if pkg.Ext == nil {
		t.Error("Expected a value for ext property")
		return
	}
}
