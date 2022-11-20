package global

type AppStatus string

const (
	Installed      AppStatus = "installed"
	UnInstalled    AppStatus = "uninstalled"
	Downloaded     AppStatus = "downloaded"
	Downloading    AppStatus = "downloading"
	Installing     AppStatus = "installing"
	Uninstalling   AppStatus = "uninstalling"
	InstallDepends AppStatus = "wait_depends" // wait for dependencies to install
	Cancelled      AppStatus = "cancelled"
)
