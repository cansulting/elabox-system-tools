//go:build windows
// +build windows

package path

// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// SystemPath.go
// Constant and variables used by the system.
// Reference: https://help.ubuntu.com/community/LinuxFilesystemTreeOverview

const PATH_SYSTEM = "c:/elabox"                      // where ela binaries will be stored
const PATH_CACHES = PATH_SYSTEM + "/caches"          // dir where caches will be saved
const PATH_HOME = PATH_SYSTEM + "/home"              // the root path for elabox. the root directory for non system apps and data
const PATH_SYSTEM_DATA = PATH_SYSTEM + "/data"       // dir where system data will be persist
const PATH_HOME_APPS = PATH_HOME + "/apps"           // where non system bin/apps will be installed
const PATH_HOME_DATA = PATH_HOME + "/data"           // where non system bin/apps data will be persist
const PATH_DOWNLOADS = PATH_HOME_DATA + "/downloads" // where downloaded files will be stored
const PATH_SYSTEM_WWW = PATH_SYSTEM + "/www"         //
const PATH_EXTERNAL_WWW = PATH_HOME + "/www"         //
const PATH_LIB = PATH_SYSTEM + "/lib"                //
const PATH_USERS = PATH_SYSTEM + "/users"            // users directory
const PATH_HOME_DOCUMENTS = PATH_HOME + "/documents"
const PATH_LOG = PATH_SYSTEM + "/log"

func GetSystemAppDir() string {
	return PATH_SYSTEM
}

// external app dir
func GetExternalAppDir() string {
	return PATH_HOME + "/apps"
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

func GetSystemAppDirData(packageId string) string {
	return PATH_SYSTEM_DATA + "/" + packageId
}

func GetExternalAppData(packageId string) string {
	return PATH_HOME_DATA + "/" + packageId
}

// get the app main executable
func GetAppInstallLocation(packageId string, external bool) string {
	if external {
		return GetExternalAppDir() + "/" + packageId
	} else {
		return GetSystemAppDir() + "/" + packageId
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
