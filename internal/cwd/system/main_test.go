package main

import (
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/servicecenter"
	"testing"
)

func TestRegistration(test *testing.T) {
	global.Initialize()
	servicecenter.Initialize(true)
	if err := appman.Initialize(true); err != nil {
		test.Error(err)
		return
	}
	_, err := appman.RegisterPackageSrc("../../builds/ela/system/apps/ela.os")
	if err != nil {
		test.Error(err)
	}
}

func TestRetrievePackage(test *testing.T) {
	global.Initialize()
	servicecenter.Initialize(true)
	if err := appman.Initialize(true); err != nil {
		test.Error(err)
		return
	}
	pk, err := appman.RetrievePackage("ela.os")
	if err != nil {
		test.Error(err)
	}
	test.Log(pk)
}
