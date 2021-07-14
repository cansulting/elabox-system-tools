package main

import (
	"archive/zip"
	"bytes"
	"ela/foundation/app/data"
	"ela/foundation/app/service"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/event"
	eventd "ela/foundation/event/data"
	"ela/foundation/path"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

/*
	installer.go
	structure for installing packages to ela system
*/
type installer struct {
	backup        *Backup             // backup instance
	BackupEnabled bool                // true if instance will create a backup for replaced files
	SilentInstall bool                // true if will install via command line false if uses actions and broadcast to others
	PackageInfo   *data.PackageConfig // package info for installer
	subinstaller  []*installer        // list of subpackages/subinstaller
}

var isSystemStopped bool = false

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

func (t *installer) initializeAppDirs() {
	dataDir := path.GetExternalAppData(t.PackageInfo.PackageId)
	if t.PackageInfo.IsSystemPackage() || !path.HasExternal() {
		dataDir = path.GetSystemAppData(t.PackageInfo.PackageId)
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
	subPackage := installer{SilentInstall: t.SilentInstall}
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
		instance.backup = &Backup{
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
	log.Println("Registering package " + t.PackageInfo.PackageId)
	// silent install means register via service
	if !t.SilentInstall {
		res, err := appController.SystemService.RequestFor(eventd.Action{})
		if err != nil {
			return errors.SystemNew("Couldnt register package ", err)
		}
		log.Println(res)
	} else {
		//register via commandline
		systemServiceSrc := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
		// check if main bin exist on path
		if _, err := os.Stat(systemServiceSrc); err != nil {
			log.Println("Registration skip for ", t.PackageInfo.PackageId, " System service not found @", systemServiceSrc)
			return nil
		}
		// run package registration
		isExternal := !t.PackageInfo.IsSystemPackage()
		installDir := path.GetExternalApp() + "/" + t.PackageInfo.PackageId
		mainExec := path.GetAppMain(t.PackageInfo.PackageId, isExternal)
		if !isExternal {
			installDir = path.GetSystemApp() + "/" + t.PackageInfo.PackageId
		}
		// check if theres a binary then register
		if _, err := os.Stat(mainExec); err == nil {
			cmd := exec.Command(systemServiceSrc, "package-reg", installDir)
			output, err := cmd.CombinedOutput()
			log.Println(string(output))
			if err != nil {
				return errors.SystemNew("installer.RegisterPackage failed", err)
			}
		} else {
			log.Println("Registration skipped for", t.PackageInfo.PackageId, "executable not found")
		}
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

// use to turn off the system
func (instance *installer) TerminateSystem() error {
	// step: connect to system
	var systemService service.IConnection
	if appController == nil || appController.Connector == nil {
		connector := event.CreateClientConnector()
		if err := connector.Open(TERMINATE_TIMEOUT); err != nil {
			return errors.SystemNew("Failed to terminate system. Unable to connect to system via connector.", err)
		}
		var err error
		systemService, err = service.NewConnection(
			connector,
			constants.SYSTEM_SERVICE_ID,
			func(message string, data interface{}) {

			})
		if err != nil {
			return errors.SystemNew("Failed to terminate system. Unable to connect to system via connector.", err)
		}
	} else {
		systemService = appController.SystemService
	}
	// step: send update mode
	response, err := systemService.RequestFor(eventd.Action{Id: constants.SYSTEM_UPDATE_MODE})
	if err != nil {
		return errors.SystemNew("Failed to terminate system. Unable to connect to system via connector.", err)
	}
	log.Println("System Will Start Update mode", response.ToString())
	return nil
}

func (instance *installer) RestartSystem() error {
	systemPath := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
	cmd := exec.Command(systemPath)
	cmd.Dir = filepath.Dir(systemPath)
	if err := cmd.Start(); err != nil {
		return errors.SystemNew("Restart system failed", err)
	}
	time.Sleep(time.Second * 3)
	os.Exit(0)
	return nil
}
