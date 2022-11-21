// this file controlls the flow of current download and install process

package installer

import (
	"errors"

	"strconv"
	"time"

	data2 "github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/broadcast"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/downloader"
	reg "github.com/cansulting/elabox-system-tools/registry/app"
)

const RATIO_VS_DOWNLOAD = 0.25 // the progress ratio of install status vs download status

// struct the handles the lifecycle of the installation
type Task struct {
	Id              string
	Url             string // this can be an http link or ipfs cid
	downloadTask    *downloader.Task
	Status          global.AppStatus
	ErrorCode       int16
	installProgress int16 // the install progress. download progress is not included
	OnStateChanged  func(task *Task)
	OnErrCallback   func(code int16, reason string)
	Dependencies    []string
	installing      bool
}

func (instance *Task) IsInstalling() bool {
	return instance.installing
}

// function that sets the current task status
func (instance *Task) setStatus(status global.AppStatus) {
	if instance.Status == status {
		return
	}
	logger.GetInstance().Debug().Msg(instance.Id + " status changed to " + string(status))
	instance.Status = status
	broadcast.PublishInstallState(instance.Id, status)
	instance.OnStateChanged(instance)
}

// function that sets install progress
func (instance *Task) SetInstallProgress(progress int16) {
	instance.installProgress = progress
	broadcast.PublishInstallProgress(uint(instance.GetOverallProgress()), instance.Id)
}

// function that returns the current download and install progress
func (instance *Task) GetOverallProgress() int16 {
	if instance.downloadTask == nil {
		return 0
	}
	// compute the overall progress
	var downloadRatio float32 = 1 - RATIO_VS_DOWNLOAD
	res := float32(instance.downloadTask.GetProgress()) * downloadRatio
	res += float32(instance.installProgress) * RATIO_VS_DOWNLOAD
	return int16(res)
}

func (instance *Task) GetDownloadPath() string {
	return instance.downloadTask.GetPath()
}

// function that starts the download
func (instance *Task) download(restart bool) {
	instance.setStatus(global.Downloading)
	if instance.downloadTask == nil {
		instance.downloadTask = downloader.AddDownload(instance.Id, instance.Url, downloader.IPFS)
		instance.downloadTask.OnStateChanged = instance.onDownloadStateChanged
		instance.downloadTask.OnProgressChanged = instance.onDownloadProgressChanged
	} else {
		if restart {
			instance.downloadTask.Reset()
		}
	}
	if err := instance.downloadTask.Start(); err != nil {
		instance.onError(global.DOWNLOAD_ERROR, err.Error())
		instance.setStatus(global.UnInstalled)
	}
}

// callback when download task state changed
func (instance *Task) onDownloadStateChanged(task *downloader.Task) {
	switch task.GetStatus() {
	case downloader.Finished:
		instance.setStatus(global.Downloaded)
	case downloader.Stopped:
		instance.setStatus(global.UnInstalled)
	case downloader.Error:
		instance.onError(task.GetError(), "download error")
		instance.setStatus(global.UnInstalled)
	}
}

// callback when download task progress changed
func (instance *Task) onDownloadProgressChanged(task *downloader.Task) {
	if err := broadcast.PublishInstallProgress(uint(instance.GetOverallProgress()), instance.Id); err != nil {
		logger.GetInstance().Error().Err(err).Caller().Msg("publish download progress failed")
	}
}

// callback when install finished
func (instance *Task) onInstalledFinished() {
	instance.installing = false
	instance.setStatus(global.Installed)
}

// callback when error found while installing
func (instance *Task) onError(code int16, reason string) {
	instance.installing = false
	logger.GetInstance().Error().Str("code", strconv.Itoa(int(code))).Caller().Msg(reason)
	instance.ErrorCode = code
	instance.SetInstallProgress(0)
	instance.OnErrCallback(code, reason)
}

// function that installs the package given the package path
// @param pkgPath the path to the package. If empty, the download path will be used
func (instance *Task) install(pkgPath string) error {
	if pkgPath == "" {
		pkgPath = instance.GetDownloadPath()
	}
	instance.setStatus(global.Installing)
	// call the package installer
	action := data.NewAction(constants.ACTION_APP_INSTALL, "", pkgPath)
	_, err := global.RPC.StartActivity(action)
	//sres, _ := res.ToSimpleResponse()
	//println(sres.Code, sres.Message)
	if err != nil {
		instance.onError(global.INSTALLER_PACKAGE_ERROR, "install error. "+err.Error())
		instance.setStatus(global.UnInstalled)
		return err
	}
	return nil
}

// function that installs the package given the package path
// @param pkgPath the path to the package. If empty, the download path will be used
func (instance *Task) Uninstall() error {
	// call the package installer
	instance.setStatus(global.Uninstalling)
	action := data.NewAction(constants.ACTION_APP_UNINSTALL, "", instance.Id)
	_, err := global.RPC.StartActivity(action)
	if err != nil {
		instance.onError(global.INSTALLER_PACKAGE_ERROR, "uninstall error "+err.Error())
		instance.setStatus(global.Installed)
		return err
	}
	return nil
}

// function use to initiate install task
func (instance *Task) Start() {
	if instance.installing {
		return
	}
	go func() {
		instance.ErrorCode = 0
		instance.installing = true
		if err := instance.waitForDependencies(); err != nil {
			return
		}
		instance.download(false)
	}()
}

// this skips the download and install the package right away given the package path
func (instance *Task) StartFromFile(pkgPath string) error {
	return instance.install(pkgPath)
}
func (instance *Task) onCancel() {
	instance.installing = false
	instance.downloadTask.Stop()
}

// callback when download task was removed from manager
func (instance *Task) onDestroy() {
	logger.GetInstance().Debug().Msg(instance.Id + " was removed from installer manager")
}

// wait for dependencies to finished
func (instance *Task) waitForDependencies() error {
	if len(instance.Dependencies) == 0 {
		return nil
	}
	logger.GetInstance().Debug().Msg("start installing dependencies")
	instance.setStatus(global.InstallDepends)
	deps := instance.Dependencies
	depTotal := len(deps)
	depInstalled := 0
	var currentDep *Task = nil
	success := false
	for {
		if currentDep == nil {
			// check if the package is installed or not
			isreg, err := reg.IsPackageInstalled(deps[0])
			if err != nil {
				return err
			}
			// install it now
			if !isreg {
				currentDep, err = CreateInstallTask(deps[0], data2.Production)
				if err != nil {
					return err
				}
				if len(deps) > 1 {
					deps = deps[1:]
				}
			} else {
				// already installed. skip package
				depInstalled++
			}
		} else {
			if currentDep.ErrorCode != 0 {
				break
			}
		}
		if currentDep != nil {
			if currentDep.Status == global.Installed {
				depInstalled++
				currentDep = nil
			} else if !currentDep.IsInstalling() {
				currentDep.Start()
			}
		}
		// did we installed all dependencies?
		if depInstalled == depTotal {
			success = true
			break
		}
		time.Sleep(time.Second)
	}
	if !success {
		instance.onError(global.INSTALL_DEPENDENCY_ERROR, "failed installing "+currentDep.Id)
		instance.setStatus(global.UnInstalled)
		return errors.New("failed installing one of the dependencies")
	}
	logger.GetInstance().Debug().Msg("finished installing dependencies")
	return nil
}
