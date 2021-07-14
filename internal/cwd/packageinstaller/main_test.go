package main

import (
	"ela/foundation/app"
	"ela/foundation/event/data"
	testapp "ela/foundation/testing/app"
	"testing"
)

// test install a package and register it
func TestSilentInstall(test *testing.T) {
	newInstall := installer{BackupEnabled: true, SilentInstall: true}
	err := newInstall.Decompress("../../builds/packages/system.ela")
	if err != nil {
		test.Error(err)
	}
	if err := newInstall.RegisterPackage(); err != nil {
		test.Error(err)
		return
	}
}

func TestRunActivityManually(t *testing.T) {
	InitializePath()
	pkgPath := "../../builds/packages/packageinstaller.ela"
	//pkgPath := `C:\Users\Jhoemar\Documents\Projects\Elabox\system-tools\internal\builds\packages`
	controller, err := app.NewController(&activity{}, nil)
	if err != nil {
		t.Error(err)
	}
	appController = controller
	testapp.RunTestApp(
		controller,
		data.ActionGroup{
			Activity: data.NewAction("", "", pkgPath)})
}
