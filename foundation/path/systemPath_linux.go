package path

/*
	SystemPath.go
	Constant and variables used by the system.
	Reference: https://help.ubuntu.com/community/LinuxFilesystemTreeOverview
*/

const PATH_SYSTEM = "/usr/ela/system"              // where ela binaries will be stored
const PATH_CACHES = "/tmp/ela"                     // dir where caches will be saved
const PATH_HOME = "/home/elabox"                   // the root path for elabox. the root directory for non system apps and data
const PATH_SYSTEM_DATA = "/var/ela/data"           // dir where system data will be persist
const PATH_APPS = PATH_HOME + "/apps"              // where non system bin/apps will be installed
const PATH_APPDATA = PATH_HOME + "/data"           // where non system bin/apps data will be persist
const PATH_DOWNLOADS = PATH_APPDATA + "/downloads" // where downloaded files will be stored
const PATH_SYSTEM_WWW = "/var/www"
const PATH_EXTERNAL_WWW = PATH_HOME + "/www"
const MAIN_EXEC_NAME = "main"
const PATH_LIB = "/usr/local/lib/ela"

func GetSystemApp() string {
	return PATH_SYSTEM
}

func GetExternalApp() string {
	return PATH_HOME
}

func GetSystemWWW() string {
	return PATH_SYSTEM_WWW
}

func GetExternalWWW() string {
	return PATH_EXTERNAL_WWW
}

// return path for system backup
func GetDefaultBackupPath() string {
	return PATH_CACHES + "/backup"
}

func GetSystemAppData(packageId string) string {
	return PATH_SYSTEM_DATA + "/" + packageId
}

func GetExternalAppData(packageId string) string {
	return PATH_APPDATA + "/" + packageId
}

// get the app main executable
func GetAppMain(packageId string, external bool) string {
	if external {
		return GetExternalApp() + "/" + packageId + "/main"
	} else {
		return GetSystemApp() + "/" + packageId + "/main"
	}
}

func GetLibPath() string {
	return PATH_LIB
}

func GetCacheDir() string {
	return PATH_CACHES
}

// return true if external is exist
func HasExternal() bool {
	return true
}
