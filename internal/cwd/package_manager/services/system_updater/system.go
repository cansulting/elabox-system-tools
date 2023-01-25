package system_updater

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	appdata "github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/broadcast"
	pkdata "github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/installer"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

type SysVer struct {
	Version string `json:"version"`
	Build   int16  `json:"build"`
}

var currentStatus global.AppStatus
var currentTask *installer.Task
var currentSysVerInfo *SysVer
var latestSysVerInfo *SysVer

const DOWNLOAD_CACHE = global.CacheDir + "/system"
const INSTALLER_CMD = "packageinstaller"

func Init() error {
	// retrieve currentSystem version
	pkinfo, err := app.RetrievePackage(constants.SYSTEM_SERVICE_ID)
	if err != nil {
		return err
	}
	currentSysVerInfo = &SysVer{
		Version: pkinfo.Version,
		Build:   pkinfo.Build,
	}
	latestSysVerInfo = &SysVer{
		Version: currentSysVerInfo.Version,
		Build:   currentSysVerInfo.Build,
	}
	// retrieve every time x for latest system version
	go func() {
		for {
			retrieveLatestSysVersion()
			time.Sleep(time.Second * global.RetrieveSystem_Delay)
		}
	}()

	return nil
}

func retrieveLatestSysVersion() {
	curVer := GetCurrentSysVersion()
	build := curVer.Build
	var lastPkInfo io.ReadCloser
	for {
		build++
		res, err := http.Get(global.SYSVER_HOST + "/" + fmt.Sprintf("%d", build) + ".json")
		if err != nil || res.StatusCode != 200 {
			if lastPkInfo == nil {
				return
			}
			// we got the latest version
			// read values
			pkconfig := appdata.DefaultPackage()
			if err := pkconfig.LoadFromReader(lastPkInfo); err != nil {
				logger.GetInstance().Error().Err(err).Msg("failed to read remote package info")
			}
			latestSysVerInfo.Build = pkconfig.Build
			latestSysVerInfo.Version = pkconfig.Version
			return
		}
		lastPkInfo = res.Body
	}
}

func GetCurrentSysVersion() *SysVer {
	return currentSysVerInfo
}

func GetInstallStatus() global.AppStatus {
	return currentStatus
}

func onStatusChanged(status global.AppStatus) {
	currentStatus = status
}

func Download(link pkdata.InstallDef) error {
	if currentStatus == global.Downloading {
		broadcast.PublishError(constants.SYSTEM_SERVICE_ID, 300, "download already in-progress")
		return errors.New("download is in progress")
	}
	task, err := installer.CreateInstallTask(link, nil, false)
	if err != nil {
		return err
	}
	task.Start()
	currentTask = task
	onStatusChanged(global.Downloading)
	return nil
}

func InstallUpdate() error {
	if currentStatus == global.Installing {
		return errors.New("update is in-progress")
	}
	if currentTask == nil {
		return errors.New("update is not yet downloaded")
	}
	logger.GetInstance().Info().Msg("Initiating system update")
	// execute package installer command
	onStatusChanged(global.Installing)
	cmd := exec.Command(INSTALLER_CMD, currentTask.GetDownloadPath(), "-s", "-l", "r")
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Process.Release(); err != nil {
		return err
	}
	// terminate elabox
	_, err := global.RPC.CallSystem(data.NewActionById(constants.SYSTEM_TERMINATE_NOW))
	return err
}
