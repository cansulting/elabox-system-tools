package path


func GetSystemAppDir() string {
	path := "C:"
	path += "\\ela\\system\\apps"
	return path
}

func GetExternalAppDir() string {
	path := "c:"
	path += "\\ela\\external\\apps"
	return path
}

func GetSystemWWW() string {
	path := "c:"
	path += "\\ela\\system\\www"
	return path
}

func GetExternalWWW() string {
	path := "c:"
	path += "\\ela\\external\\www"
	return path
}

// return path for system backup
func GetDefaultBackupPath() string {
	path := "c:"
	path += "\\ela\\system\\backup"
	return path
}

func GetSystemAppDirData(packageId string) string {
	path := "c:"
	path += "\\ela\\system\\data\\" + packageId
	return path
}

func GetExternalAppData(packageId string) string {
	path := "c:"
	path += "\\ela\\external\\data\\" + packageId
	return path
}

// get the app main executable
func GetAppMain(packageId string, external bool) string {
	if external {
		return GetExternalAppDir() + "\\" + packageId + "\\" + 
	} else {
		return GetSystemAppDir() + "\\" + packageId + "\\" + 
	}
}

func GetCacheDir() string {
	path := "c:"
	path += "\\ela\\caches"
	return path
}

func GetLibPath() string {
	path := "c:"
	path += "\\ela\\lib"
	return path
}

// return true if external is exist
func HasExternal() bool {
	return true
}
