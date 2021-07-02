package path

import (
	"os"
)

func GetSystemApp() string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/system/apps"
	return path
}

func GetExternalApp() string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/external/apps"
	return path
}

func GetSystemWWW() string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/system/www"
	return path
}

func GetExternalWWW() string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/external/www"
	return path
}

// return path for system backup
func GetDefaultBackupPath() string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/system/backup"
	return path
}

func GetSystemAppData(packageId string) string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/system/data/" + packageId
	return path
}

func GetExternalAppData(packageId string) string {
	path, _ := os.Getwd()
	path += "/../../builds"
	path += "/ela/external/data/" + packageId
	return path
}

// get the app main executable
func GetAppMain(packageId string, external bool) string {
	if external {
		return GetExternalApp() + "/" + packageId + "/main.exe"
	} else {
		return GetSystemApp() + "/" + packageId + "/main.exe"
	}
}

// return true if external is exist
func HasExternal() bool {
	return true
}
