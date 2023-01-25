//go:build RELEASE || STAGING
// +build RELEASE STAGING

package system_updater

func GetLatestSysVersion() *SysVer {
	return latestSysVerInfo
}

func DownloadLatest() error {
	latest := GetLatestSysVersion()
	link := pkdata.InstallDef{
		Url:  global.SYSVER_HOST + "/" + fmt.Sprintf("%d", latest.Build) + ".box",
		Name: "Elabox",
		Id:   constants.SYSTEM_SERVICE_ID,
	}
	return Download(link)
}
