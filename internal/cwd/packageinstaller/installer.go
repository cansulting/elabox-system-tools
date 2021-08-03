package main

import (
	"archive/zip"
	"bytes"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/foundation/perm"
	"ela/internal/cwd/global"
	"ela/internal/cwd/packageinstaller/pkg"
	"ela/internal/cwd/packageinstaller/utils"
	"ela/registry/app"
	"io"
	"log"
	"os"
)

/*
	installer.go
	structure for installing packages to ela system
*/
type installer struct {
	backup         *utils.Backup       // backup instance
	BackupEnabled  bool                // true if instance will create a backup for replaced files
	packageInfo    *data.PackageConfig // package info for installer
	packageContent *pkg.Data
	subinstaller   []*installer // list of subpackages/subinstaller
	onProgress     func(uint16)
	progress       uint16
}

func NewInstaller(content *pkg.Data, backup bool) *installer {
	return &installer{
		BackupEnabled:  backup,
		packageInfo:    content.Config,
		packageContent: content,
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
		backupPath := path.GetDefaultBackupPath() + "/system.backup"
		err := instance.backup.Create(backupPath)
		if err != nil {
			return errors.SystemNew("Couldn't create backup for @"+backupPath, err)
		}
	}
	// step: create backup for app bin
	if instance.BackupEnabled {
		if err := instance.backup.AddFiles(packageInfo.GetInstallDir()); err != nil {
			log.Println("installer failed to backup app dir "+packageInfo.GetInstallDir(), err.Error(), "continue...")
		}
	}
	// preinstall
	if err := instance.preinstall(); err != nil {
		return err
	}
	// step: delete the package first
	if err := utils.UninstallPackage(packageInfo.PackageId, false); err != nil {
		return err
	}
	log.Println("installer:start installing ", packageInfo.PackageId)
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
	// step: iterate each file and save it
	for _, file := range instance.packageContent.Files {
		// step: open source file
		log.Println("installer:decompress() extracting", file.Name)
		reader, err := file.Open()
		if err != nil {
			log.Println("installer::decompress error", file.Name, err)
			return err
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
				log.Println("error", "installer::uncompress to file ", file.Name, "...", err)
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
			log.Println("installer no Filter. skipped ", targetPath)
			continue
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
		log.Println("installer.initializeAppDirs failed", err)
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
	subPackage := NewInstaller(pkg, false)
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
	log.Println("Finalizing installer")
	if err := t.registerPackage(); err != nil {
		return errors.SystemNew("Unable to register package "+t.backup.PackageId, err)
	}
	if t.packageContent.HasPostInstallScript() {
		if err := t.packageContent.StartPostInstall(); err != nil {
			return err
		}
	}
	t.packageContent.Clean()
	return nil
}

// use to revert changes based from backup
func (t *installer) RevertChanges() error {
	if t.BackupEnabled && t.backup != nil {
		log.Println("Reverting changes...")
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
			log.Println(err.Error())
		}
	}
}

func (instance *installer) SetProgressListener(listener func(uint16)) {
	instance.onProgress = listener
}

func (instance *installer) setProgress(progress uint16, log string) {
	instance.progress = progress
	if instance.onProgress != nil {
		instance.onProgress(progress)
	}
}

func (instance *installer) GetProgress() uint16 {
	return instance.progress
}
