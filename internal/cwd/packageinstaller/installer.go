package main

import (
	"archive/zip"
	"bytes"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
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
	backup        *utils.Backup       // backup instance
	BackupEnabled bool                // true if instance will create a backup for replaced files
	PackageInfo   *data.PackageConfig // package info for installer
	subinstaller  []*installer        // list of subpackages/subinstaller
}

// use to uncompress the file to target
func (instance *installer) Decompress(sourceFile string) error {
	// step: read package
	z, err := zip.OpenReader(sourceFile)
	if err != nil {
		return errors.SystemNew("installer.Decompress() failed to locate "+sourceFile, err)
	}
	defer z.Close()
	return instance.decompressFromReader(z.File)
}

// decompress package based from reader
func (instance *installer) decompressFromReader(files []*zip.File) error {
	// step: load package
	packageInfo := instance.PackageInfo
	if packageInfo == nil {
		packageInfo = data.DefaultPackage()
		if err := packageInfo.LoadFromZipFiles(files); err != nil {
			return errors.SystemNew("installer.decompressFromReader unable to load package info", err)
		}
		instance.PackageInfo = packageInfo
	}
	if !instance.PackageInfo.IsValid() {
		return errors.SystemNew("installer.decompressFromReader invalid package info", nil)
	}
	// step: create backup for app bin
	if err := instance.backup.AddFiles(packageInfo.GetInstallDir()); err != nil {
		log.Println("installer failed to backup app dir "+packageInfo.GetInstallDir(), "continue...")
	}
	// step: delete the package first
	if err := utils.UninstallPackage(packageInfo.PackageId, false); err != nil {
		return err
	}
	log.Println("installer:start installing ", packageInfo.PackageId, "silent=")
	// step: init install location and filters
	appInstallPath, wwwInstallPath := _getInstallLocation(packageInfo)
	filters := []filter{
		// bin
		{keyword: "bin", rename: packageInfo.PackageId, installTo: appInstallPath},
		// library
		{keyword: "lib", rename: packageInfo.PackageId, installTo: path.GetLibPath()},
		// www
		{keyword: "www", rename: packageInfo.PackageId, installTo: wwwInstallPath},
		{keyword: constants.APP_CONFIG_NAME, installTo: appInstallPath + "/" + packageInfo.PackageId},
		// subpackages
		{keyword: "packages/", customProcess: instance._onSubPackage},
		// node js
		{keyword: "nodejs", installTo: packageInfo.GetInstallDir()},
	}
	// step: iterate each file and save it
	for _, file := range files {
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
		//step: apply filter and resolve directories
		filterApplied := false
		for _, filter := range filters {
			// use filter to customize the destination or change name
			newPath, err, applied := filter.applyTo(targetPath, 0764, reader, file.CompressedSize64)
			if err != nil {
				log.Println(err)
				return nil
			}
			// filter was applied. break
			if applied {
				targetPath = newPath
				filterApplied = true
				break
			}
		}
		// no filter was applied. use the default destination
		if !filterApplied {
			log.Println("installer no filter. skipped ", targetPath)
			continue
		}
		// step: is valid target path?
		if targetPath != "" {
			// step: check if instance file already exist. then create backup
			if instance.BackupEnabled {
				if os.Stat(targetPath); err == nil {
					err = instance.createBackupFor(targetPath)
					if err != nil {
						return err
					}
				}
			}
			// step: create dest file
			newFile, err := os.Create(targetPath)
			if err != nil {
				log.Println("error", "installer::uncompress to file ", file.Name, "...", err)
				return err
			}
			// step: write to file
			io.Copy(newFile, reader)
			newFile.Close()
		}
	}
	instance.initializeAppDirs()
	instance._closeBackup()
	return nil
}

// initialize directories
func (t *installer) initializeAppDirs() {
	dataDir := path.GetExternalAppData(t.PackageInfo.PackageId)
	if t.PackageInfo.IsSystemPackage() || !path.HasExternal() {
		dataDir = path.GetSystemAppData(t.PackageInfo.PackageId)
		t.PackageInfo.ChangeToSystemLocation()
	}
	if err := os.MkdirAll(dataDir, 0740); err != nil {
		log.Println("installer.initializeAppDirs failed", err)
	}
}

// return app and www install path base on the package
func _getInstallLocation(packageInfo *data.PackageConfig) (string, string) {
	appInstallPath := path.GetExternalApp()
	wwwInstallPath := path.GetExternalWWW()
	if packageInfo.IsSystemPackage() ||
		!path.HasExternal() {
		appInstallPath = path.GetSystemApp()
		wwwInstallPath = path.GetSystemWWW()
	}
	return appInstallPath, wwwInstallPath
}

// callback when theres a subpackage
func (t *installer) _onSubPackage(path string, reader io.ReadCloser, size uint64) error {
	subPackage := installer{}
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
	// step: decompress subpackage file
	if err := subPackage.decompressFromReader(newReader.File); err != nil {
		return errors.SystemNew("installer: subpackage "+path+"...", err)
	}
	if t.subinstaller == nil {
		t.subinstaller = make([]*installer, 0, 4)
	}
	// step: add to list
	t.subinstaller = append(t.subinstaller, &subPackage)
	return nil
}

// create backup for file
func (instance *installer) createBackupFor(src string) error {
	if instance.backup == nil {
		instance.backup = &utils.Backup{
			PackageId: "",
		}
		backupPath := path.GetDefaultBackupPath() + "/system.backup"
		err := instance.backup.Create(backupPath)
		if err != nil {
			return errors.SystemNew("Couldn't create backup for "+src+".", err)
		}
	}
	instance.backup.AddFile(src)
	return nil
}

// start registering the package and sub packages
func (t *installer) RegisterPackage() error {
	if err := app.RegisterPackage(t.PackageInfo); err != nil {
		return err
	}
	// register subpackages
	for _, subinstall := range t.subinstaller {
		if err := subinstall.RegisterPackage(); err != nil {
			return err
		}
	}
	return nil
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
