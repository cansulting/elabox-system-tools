package main

import (
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/servicecenter"
	"testing"
)

func TestPacakageRegistration(test *testing.T) {
	//global.Initialize()
	servicecenter.Initialize(true)
	if err := appman.Initialize(true); err != nil {
		test.Error(err)
		return
	}
	_, err := appman.RegisterPackageSrc("../../builds/ela/system/apps/ela.system")
	if err != nil {
		test.Error(err)
	}
}

func TestRetrievePackage(test *testing.T) {
	//global.Initialize()
	servicecenter.Initialize(true)
	if err := appman.Initialize(true); err != nil {
		test.Error(err)
		return
	}
	pk, err := appman.RetrievePackage("ela.installer")
	if err != nil {
		test.Error(err)
	}
	test.Log(pk)
}
