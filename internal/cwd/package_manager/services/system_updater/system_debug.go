//go:build !RELEASE && !STAGING
// +build !RELEASE,!STAGING

package system_updater

import (
	"fmt"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	pkdata "github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
)

func GetLatestSysVersion() *SysVer {
	return &SysVer{
		Build:   latestSysVerInfo.Build + 1,
		Version: "samplever",
	}
}

func DownloadLatest() error {
	cur := GetCurrentSysVersion()
	link := pkdata.InstallDef{
		Url:  global.SYSVER_HOST + "/" + fmt.Sprintf("%d", cur.Build) + ".box",
		Name: "Elabox",
		Id:   constants.SYSTEM_SERVICE_ID,
	}
	return Download(link)
}
