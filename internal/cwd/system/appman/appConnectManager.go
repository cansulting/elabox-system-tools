package appman

// currently running processes
var running map[string]*AppConnect

// use to run process for specific package. return true if success, false if already running
func getAppConnect(packageId string) *AppConnect {
	app, ok := running[packageId]
	// is already running? return false
	if ok {
		return app
	}
	location, _ := retrievePackageSource(packageId)
	if location == "" {
		return nil
	}
	app = &AppConnect{
		packageId: packageId, location: location}
	running[packageId] = app
	return app
}

func lookupAppConnect(packageId string) *AppConnect {
	pk, ok := running[packageId]
	// is already running? return false
	if ok {
		return pk
	}
	return nil
}

func removeAppConnect(packageId string, terminate bool) {
	app := lookupAppConnect(packageId)
	if app != nil {
		if terminate {
			app.forceTerminate()
		}
		delete(running, packageId)
	}
}
