package constants

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
