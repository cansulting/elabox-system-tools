package main

import (
	"archive/zip"
	"bytes"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	eventd "ela/foundation/event/data"
	"ela/foundation/path"
	"io"
	"log"
	"os"
	"os/exec"
)

/*
	installer.go
	structure for installing packages to ela system
*/
type installer struct {
	backup        *Backup             // backup instance
	BackupEnabled bool                // true if instance will create a backup for replaced files
	SilentInstall bool                // true if will install via command line false if uses actions and broadcast to others
	packageInfo   *data.PackageConfig // package info for installer
	subinstaller  []*installer        // list of subpackages/subinstaller
}

const WRITE_SIZE = 1000

// use to uncompress the file to target
func (instance *installer) decompress(sourceFile string) error {
	// step: read package
	z, err := zip.OpenReader(sourceFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer z.Close()
	return instance.decompressFromReader(z.File)
}

// decompress package based from reader
func (instance *installer) decompressFromReader(files []*zip.File) error {
	// step: load package
	packageInfo, error := instance._loadPackage(files)
	instance.packageInfo = packageInfo
	if error != nil {
		return error
	}
	log.Println("installer:start installing ", packageInfo.PackageId, "silent=", instance.SilentInstall)
	// step: init install location and filters
	appInstallPath, wwwInstallPath := _getInstallLocation(packageInfo)
	filters := []filter{
		// skip
		{keyword: "*.DS_Store"},
		// move to location
		{keyword: "bin", rename: packageInfo.PackageId, installTo: appInstallPath},
		{keyword: "www", rename: packageInfo.PackageId, installTo: wwwInstallPath},
		{keyword: constants.APP_CONFIG_NAME, installTo: appInstallPath + "/" + packageInfo.PackageId},
		// subpackage
		{keyword: "packages/", customProcess: instance._onSubPackage},
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
					error = instance.createBackupFor(targetPath)
					if error != nil {
						return error
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
		}
	}
	instance._closeBackup()
	if err := packageInfo.GetError(); err != nil {
		return err
	}
	return nil
}

// return app and www install path base on the package
func _getInstallLocation(packageInfo *data.PackageConfig) (string, string) {
	appInstallPath := path.GetExternalApp()
	wwwInstallPath := path.GetExternalWWW()
	if packageInfo.InstallLocation == "system" ||
		!path.HasExternal() {
		appInstallPath = path.GetSystemApp()
		wwwInstallPath = path.GetSystemWWW()
	}
	return appInstallPath, wwwInstallPath
}

// package loading
func (instance *installer) _loadPackage(files []*zip.File) (*data.PackageConfig, error) {
	packageInfo := data.DefaultPackage()
	for _, file := range files {
		if file.Name != constants.APP_CONFIG_NAME {
			continue
		}
		reader, error := file.Open()
		if error != nil {
			return nil, &InstallError{errorString: "Load Package error. " + error.Error()}
		}
		error = packageInfo.LoadFromReader(reader)
		if error != nil {
			return nil, &InstallError{errorString: "Load Package error. " + error.Error()}
		}
		break
	}
	return packageInfo, nil
}

// callback when theres a subpackage
func (t *installer) _onSubPackage(path string, reader io.ReadCloser, size uint64) error {
	subPackage := installer{SilentInstall: t.SilentInstall}
	// step: convert buffer to zip reader
	newBuffer := bytes.NewBuffer([]byte{})
	written, err := io.Copy(newBuffer, reader)
	if err != nil {
		return &InstallError{errorString: "Subpackage error " + path + "..." + err.Error()}
	}
	newReader, err := zip.NewReader(bytes.NewReader(newBuffer.Bytes()), written)
	if err != nil {
		return &InstallError{errorString: "installer: subpackage " + path + "..." + err.Error()}
	}
	// step: decompress subpackage file
	if err := subPackage.decompressFromReader(newReader.File); err != nil {
		return &InstallError{errorString: "installer: subpackage " + path + "..." + err.Error()}
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
		instance.backup = &Backup{
			PackageId: "",
		}
		backupPath := path.GetDefaultBackupPath() + "/system.backup"
		err := instance.backup.Create(backupPath)
		if err != nil {
			return &InstallError{errorString: "Couldn't create backup for " + src + "." + err.Error()}
		}
	}
	instance.backup.AddFile(src)
	return nil
}

// start registering the package and sub packages
func (t *installer) registerPackage() error {
	log.Println("Registering package " + t.packageInfo.PackageId)
	// silent install means register via service
	if !t.SilentInstall {
		res, err := appController.SystemService.RequestFor(eventd.Action{})
		if err != nil {
			return &InstallError{errorString: "Couldnt register package " + err.Error()}
		}
		log.Println(res)
	} else {
		//register via commandline
		systemServiceSrc := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
		// check if main bin exist on path
		if _, err := os.Stat(systemServiceSrc); err != nil {
			log.Println("Registration skip for ", t.packageInfo.PackageId, " System service not found @", systemServiceSrc)
			return nil
		}
		// run package registration
		installDir := path.GetExternalApp() + "/" + t.packageInfo.PackageId
		if t.packageInfo.IsSystemPackage() {
			installDir = path.GetSystemApp() + "/" + t.packageInfo.PackageId
		}
		cmd := exec.Command(systemServiceSrc, "package-reg", installDir)
		output, err := cmd.CombinedOutput()
		log.Println(string(output))
		if err != nil {
			return &InstallError{errorString: "Couldnt register package " + err.Error()}
		}
		return nil
	}
	// register subpackages
	for _, subinstall := range t.subinstaller {
		if err := subinstall.registerPackage(); err != nil {
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
