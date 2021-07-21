package main

import (
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/servicecenter"
	"ela/registry/app"
	"log"
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

func TestRetrieveAllPackages(test *testing.T) {
	//global.Initialize()
	servicecenter.Initialize(true)
	if err := appman.Initialize(true); err != nil {
		test.Error(err)
		return
	}
	pks, err := app.RetrieveAllPackages()
	if err != nil {
		test.Error(err)
	}
	for _, pk := range pks {
		log.Println(pk.ToString())
	}

}
