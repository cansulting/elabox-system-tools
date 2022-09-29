package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/internal/cwd/global"
	pkconst "github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/pkg"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/utils"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

/*
	installer.go
	structure for installing packages to ela system
*/
type installer struct {
	backup            *utils.Backup       // backup instance
	BackupEnabled     bool                // true if instance will create a backup for replaced files
	packageInfo       *data.PackageConfig // package info for installer
	packageContent    *pkg.Data
	subinstaller      []*installer // list of subpackages/subinstaller
	onProgress        func(int, string)
	onError           func(string, int, string, error)
	progress          int // current progress 0 - 100
	broadcastProgress bool
}

// create new installer instance
// @param content - package content
// @param broadcast - true if will broadcast to system for progress and status updates
func NewInstaller(content *pkg.Data, backup bool, broadcastProgress bool) *installer {
	return &installer{
		BackupEnabled:     backup,
		packageInfo:       content.Config,
		packageContent:    content,
		broadcastProgress: broadcastProgress,
	}
}

// decompress package based from reader
func (instance *installer) Start() error {
	packageInfo := instance.packageInfo
	// initialize backup
	if instance.BackupEnabled && instance.backup == nil {
		instance.backup = &utils.Backup{
			PackageId: packageInfo.PackageId,
		}
		if err := os.MkdirAll(path.GetDefaultBackupPath(), perm.PUBLIC_WRITE); err != nil {
			return errors.SystemNew("Failed to create backup directory", err)
		}
		backupPath := path.GetDefaultBackupPath() + "/system.backup"
		if err := instance.backup.Create(backupPath); err != nil {
			return errors.SystemNew("Couldn't create backup for @"+backupPath, err)
		}
	}
	// step: create backup for app bin
	if instance.BackupEnabled {
		if err := instance.backup.AddFiles(packageInfo.GetInstallDir()); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("installer failed to backup app dir " + packageInfo.GetInstallDir())
		}
	}
	// step: resolve dir for data
	os.MkdirAll(packageInfo.GetDataDir(), perm.PUBLIC_WRITE)
	// preinstall
	if err := instance.preinstall(); err != nil {
		return err
	}
	// step: delete the package first
	if err := utils.UninstallPackage(packageInfo.PackageId, false, false, instance.broadcastProgress); err != nil {
		pkconst.Logger.Error().Err(err).Caller().Msg("unable to uninstall package " + packageInfo.PackageId)
	}
	pkconst.Logger.Info().Msg("start installing " + packageInfo.PackageId)
	// step: init install location and filters
	appInstallPath, wwwInstallPath := _getInstallLocation(packageInfo)
	filters := []utils.Filter{
		// bin
		{Keyword: "bin", Rename: packageInfo.PackageId, InstallTo: appInstallPath, Perm: perm.PRIVATE},
		// library
		{Keyword: "lib", Rename: packageInfo.PackageId, InstallTo: path.GetLibPath(), Perm: perm.PUBLIC_VIEW},
		// www
		{Keyword: global.PKEY_WWW, Rename: packageInfo.PackageId, InstallTo: wwwInstallPath, Perm: perm.PUBLIC_VIEW},
		{Keyword: constants.APP_CONFIG_NAME, InstallTo: appInstallPath + "/" + packageInfo.PackageId, Perm: perm.PUBLIC_VIEW},
		// subpackages
		{Keyword: "packages/", CustomProcess: instance._onSubPackage, Perm: perm.PUBLIC},
		// node js
		{Keyword: "nodejs", InstallTo: packageInfo.GetInstallDir(), Perm: perm.PRIVATE},
	}
	fileCount := len(instance.packageContent.Files)
	processed := 0
	// step: iterate each file and save it
	for _, file := range instance.packageContent.Files {
		// step: open source file
		pkconst.Logger.Debug().Msg("extracting " + file.Name)
		reader, err := file.Open()
		if err != nil {
			return errors.SystemNew("File open error "+file.Name, err)
		}
		defer reader.Close()
		file.DataOffset()
		isDir := file.FileInfo().IsDir()
		if isDir {
			continue
		}
		targetPath := file.Name
		//step: apply Filter and resolve directories
		var filterApplied *utils.Filter
		for _, Filter := range filters {
			// use Filter to customize the destination or change name
			newPath, err, applied := Filter.CanApply(targetPath, reader, file.CompressedSize64)
			if err != nil {
				pkconst.Logger.Error().Err(err).Caller().Msg("uncompressing file " + file.Name + ". continue...")
				return nil
			}
			// Filter was applied. break
			if applied {
				filterApplied = &Filter
				if newPath != "" {
					// step: check if instance file already exist. then create backup
					if instance.BackupEnabled {
						if os.Stat(targetPath); err == nil {
							err = instance.createBackupFor(targetPath)
							if err != nil {
								return err
							}
						}
					}
					if err := Filter.Save(newPath, reader); err != nil {
						return errors.SystemNew("Unable to save "+file.Name, err)
					}
				}
				break
			}
		}
		// no Filter was applied. use the default destination
		if filterApplied == nil {
			pkconst.Logger.Warn().Msg("keyword no filter. skipped " + targetPath)
			continue
		}
		// step: compute progress
		if instance.onProgress != nil {
			processed++
			// comppute progress
			instance.setProgress(int(float32(processed) / float32(fileCount) * 100))
		}
	}
	instance.initializeAppDirs()
	instance._closeBackup()
	return nil
}

// initialize directories
func (t *installer) initializeAppDirs() {
	dataDir := path.GetExternalAppData(t.packageInfo.PackageId)
	if t.packageInfo.IsSystemPackage() || !path.HasExternal() {
		dataDir = path.GetSystemAppDirData(t.packageInfo.PackageId)
		t.packageInfo.ChangeToSystemLocation()
	}
	if err := os.MkdirAll(dataDir, perm.PUBLIC_WRITE); err != nil {
		pkconst.Logger.Error().Err(err).Caller().Msg("Mkdirall failed " + dataDir)
	}
}

// return app and www install path base on the package
func _getInstallLocation(packageInfo *data.PackageConfig) (string, string) {
	appInstallPath := path.GetExternalAppDir()
	wwwInstallPath := path.GetExternalWWW()
	if packageInfo.IsSystemPackage() ||
		!path.HasExternal() {
		appInstallPath = path.GetSystemAppDir()
		wwwInstallPath = path.GetSystemWWW()
	}
	return appInstallPath, wwwInstallPath
}

// callback when theres a subpackage
func (t *installer) _onSubPackage(path string, reader io.ReadCloser, size uint64) error {
	// step: convert buffer to zip reader
	newBuffer := bytes.NewBuffer([]byte{})
	written, err := io.Copy(newBuffer, reader)
	if err != nil {
		return errors.SystemNew("Subpackage error "+path+"...", err)
	}
	newReader, err := zip.NewReader(bytes.NewReader(newBuffer.Bytes()), written)
	if err != nil {
		return errors.SystemNew("installer: subpackage "+path+"...", err)
	}
	// initialize package
	pkg, err := pkg.LoadFromZipFiles(newReader.File)
	if err != nil {
		return errors.SystemNew("Failed loading subpackage ", err)
	}
	subPackage := NewInstaller(pkg, false, t.broadcastProgress)
	// step: decompress subpackage file
	if err := subPackage.Start(); err != nil {
		return errors.SystemNew("installer: subpackage "+path+"...", err)
	}
	if t.subinstaller == nil {
		t.subinstaller = make([]*installer, 0, 4)
	}
	// step: add to list
	t.subinstaller = append(t.subinstaller, subPackage)
	return nil
}

// create backup for file
func (instance *installer) createBackupFor(src string) error {
	instance.backup.AddFile(src)
	return nil
}

func (t *installer) preinstall() error {
	if err := t.packageContent.ExtractScripts(); err != nil {
		return errors.SystemNew("Failed to extract scripts", err)
	}
	if t.packageContent.HasPreInstallScript() {
		if err := t.packageContent.StartPreInstall(); err != nil {
			return err
		}
	}
	return nil
}

// several steps before installation finalizes
func (t *installer) Finalize() error {
	pkconst.Logger.Debug().Msg("Finalizing installer")
	if err := t.registerPackage(); err != nil {
		return errors.SystemNew("Unable to register package "+t.backup.PackageId, err)
	}
	if t.packageContent.HasPostInstallScript() {
		if err := t.packageContent.StartPostInstall(); err != nil {
			return err
		}
	}
	// activate port for activity custom port
	if t.packageInfo.ActivityGroup.CustomPort > 0 {
		if err := utils.AllowPort(t.packageInfo.ActivityGroup.CustomPort); err != nil {
			pkconst.Logger.Error().Err(err).Caller().Msg("failed to allow port " + strconv.Itoa(t.packageInfo.ActivityGroup.CustomPort) + " for " + t.packageInfo.PackageId)
		}
	}
	// activate ports
	if len(t.packageInfo.ExposePorts) > 0 {
		for _, port := range t.packageInfo.ExposePorts {
			if err := utils.AllowPort(port); err != nil {
				pkconst.Logger.Error().Err(err).Caller().Msg("failed to allow port " + strconv.Itoa(port) + " for " + t.packageInfo.PackageId)
			}
		}
	}
	t.packageContent.Clean()
	return nil
}

// use to revert changes based from backup
func (t *installer) RevertChanges() error {
	if t.BackupEnabled && t.backup != nil {
		pkconst.Logger.Debug().Msg("Reverting changes...")
		bkSrc := t.backup.GetSource()
		t._closeBackup()
		bk := utils.Backup{}
		if err := bk.LoadAndApply(bkSrc); err != nil {
			return err
		}
	}
	return nil
}

// start registering the package and sub packages
func (t *installer) registerPackage() error {
	if err := app.RegisterPackage(t.packageInfo); err != nil {
		return err
	}
	// register subpackages
	for _, subinstall := range t.subinstaller {
		if err := subinstall.registerPackage(); err != nil {
			return errors.SystemNew("Failed to register subpackage", err)
		}
	}
	return nil
}

// use to check if theres a custom installer inside zip
func (instance *installer) hasCustomInstallerInZip(files []*zip.File) bool {
	for _, file := range files {
		if file.Name == global.PACKAGEKEY_CUSTOM_INSTALLER {
			return true
		}
	}
	return false
}

// use to close the backup
func (instance *installer) _closeBackup() {
	if instance.backup != nil {
		err := instance.backup.Close()
		if err != nil {
			//instance.onError()
			pkconst.Logger.Error().Err(err).Caller().Msg("Failed to close backup.")
		}
	}
}

func (instance *installer) SetProgressListener(listener func(int, string)) {
	instance.onProgress = listener
}

// sets the listener for error callback
func (instance *installer) SetErrorListener(listener func(string, int, string, error)) {
	instance.onError = listener
}

func (instance *installer) setProgress(progress int) {
	instance.progress = progress
	if instance.onProgress != nil {
		instance.onProgress(progress, instance.packageInfo.PackageId)
	}
}

func (instance *installer) _onError(code int, reason string, err error) {
	pkconst.Logger.Error().Err(err).Caller().Msg(reason)
	if instance.onError != nil {
		instance.onError(instance.packageInfo.PackageId, code, reason, err)
	}
}

// func (instance *installer) GetProgress() uint16 {
// 	return instance.progress
// }
