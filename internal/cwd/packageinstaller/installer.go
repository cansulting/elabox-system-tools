package main

//https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool
//https://golang.org/pkg/archive/zip/#example_Reader
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
	spath "path"
)

type installer struct {
	backup        *Backup // backup instance
	backupEnabled bool    // true if this will create a backup for replaced files
	silentInstall bool    // true if will install via command line false if uses actions
}

const WRITE_SIZE = 1000

// use to initialize.
// @silentInstall true if will install via service center command line
func (t *installer) initInstall(silentInstall bool) {
	t.silentInstall = silentInstall
}

// use to uncompress the file to target
func (this *installer) decompress(sourceFile string) error {
	// step: read package
	z, err := zip.OpenReader(sourceFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer z.Close()
	return this.decompressFromReader(z.File)
}

// decompress package based from reader
func (this *installer) decompressFromReader(files []*zip.File) error {
	// step: load package
	packageInfo, error := this._loadPackage(files)
	if error != nil {
		return error
	}
	log.Println("installer:start installing ", packageInfo.PackageId, "silent=", this.silentInstall)
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
		{keyword: "packages/", customProcess: this._onSubPackage},
	}
	// step: iterate each file and save it
	for _, file := range files {
		// step: open source file
		log.Println("installer:uncompress extract start", file.Name)
		reader, err := file.Open()
		defer reader.Close()
		if err != nil {
			log.Println("installer::uncompress error", file.Name, err)
			return err
		}
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
			// step: check if this file already exist. then create backup
			if this.backupEnabled {
				if os.Stat(targetPath); err == nil {
					error = this.createBackupFor(targetPath)
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
	this.closeBackup()
	if err := packageInfo.GetError(); err != nil {
		return err
	}
	if err := this.registerPackage(packageInfo); err != nil {
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
func (this *installer) _loadPackage(files []*zip.File) (*data.PackageConfig, error) {
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
	subPackage := installer{}
	subPackage.initInstall(t.silentInstall)
	newBuffer := bytes.NewBuffer([]byte{})
	written, err := io.Copy(newBuffer, reader)
	if err != nil {
		return &InstallError{errorString: "Subpackage error " + path + "..." + err.Error()}
	}
	newReader, err := zip.NewReader(bytes.NewReader(newBuffer.Bytes()), written)
	if err != nil {
		return &InstallError{errorString: "installer: subpackage " + path + "..." + err.Error()}
	}
	if err := subPackage.decompressFromReader(newReader.File); err != nil {
		return &InstallError{errorString: "installer: subpackage " + path + "..." + err.Error()}
	}
	return nil
}

// create backup for file
func (this *installer) createBackupFor(src string) error {
	if this.backup == nil {
		this.backup = &Backup{
			PackageId: "",
		}
		backupPath := path.GetDefaultBackupPath() + "/system.backup"
		err := this.backup.Create(backupPath)
		if err != nil {
			return &InstallError{errorString: "Couldn't create backup for " + src + "." + err.Error()}
		}
	}
	this.backup.AddFile(src)
	return nil
}

// start registering the package
func (t *installer) registerPackage(pk *data.PackageConfig) error {
	log.Println("Registering package " + pk.PackageId)
	// register via service
	if !t.silentInstall {
		res, err := appController.SystemService.RequestFor(eventd.Action{})
		if err != nil {
			return &InstallError{errorString: "Couldnt register package " + err.Error()}
		}
		log.Println(res)
		return nil
	} else {
		//register via commandline
		systemServiceSrc := path.GetAppMain(constants.SYSTEM_SERVICE_ID, false)
		if _, err := os.Stat(systemServiceSrc); err != nil {
			log.Println("Registration skip for ", pk.PackageId, " System service not found @", systemServiceSrc)
			return nil
		}
		installDir := spath.Dir(pk.Source)
		cmd := exec.Command(systemServiceSrc, "package-reg "+installDir)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
}

// use to close the backup
func (this *installer) closeBackup() {
	if this.backup != nil {
		err := this.backup.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
}
