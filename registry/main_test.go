package main

import (
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/servicecenter"
	"ela/registry/app"
	"testing"
)

func TestPacakageRegistration(test *testing.T) {
	//global.Initialize()
	servicecenter.Initialize(true)
	if err := appman.Initialize(true); err != nil {
		test.Error(err)
		return
	}
	_, err := app.RegisterPackageSrc(`C:\ela\external\apps\ela.sample2`)
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
	pk, err := app.RetrievePackage("ela.sample")
	if err != nil {
		test.Error(err)
	}
	test.Log(pk)
}
