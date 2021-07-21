package utils

import (
	"archive/zip"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/foundation/perm"
	"log"
	"os"
	"os/exec"
)

/*
	CustomExec.go
	this struct handles execution of scripts and custom installer attach to it
*/

const sh = "/bin/sh"
const customInstallerName = "scripts/installer"

type CustomExec struct {
	config                 *data.PackageConfig
	preInstall             string
	postInstall            string
	IsCustomInstallerAvail bool // true if custom installer availbable
}

func getLocation() string {
	return path.GetCacheDir() + "/custominstall"
}

func getCustomInstallerLocation() string {
	return getLocation() + "/" + customInstallerName
}

func GetSHConfig(config *data.PackageConfig, files []*zip.File) (*CustomExec, error) {
	cinstallerDir := getLocation()

	// export installer and scripts
	filters := []Filter{
		{Keyword: "scripts/" + constants.PREINSTALL_SH, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
		{Keyword: "scripts/" + constants.POSTINSTALL_SH, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
		{Keyword: customInstallerName, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
	}
	found := false
	// initialize scripts and files and move to cache directory
	for _, file := range files {
		for _, cfilter := range filters {
			newPath, err, applied := cfilter.CanApply(file.Name, nil, 0)
			if err != nil {
				return nil, err
			}
			if applied {
				reader, err := file.Open()
				if err != nil {
					return nil, err
				}
				if cfilter.Save(newPath, reader); err != nil {
					return nil, errors.SystemNew("Failed saving "+file.Name, err)
				}
				reader.Close()
				found = true
			}
		}
	}
	if found {
		isInstallerExisted := false
		if _, err := os.Stat(getCustomInstallerLocation()); err == nil {
			isInstallerExisted = true
		}
		newInstance := &CustomExec{
			preInstall:             cinstallerDir + "/" + constants.PREINSTALL_SH,
			postInstall:            cinstallerDir + "/" + constants.POSTINSTALL_SH,
			IsCustomInstallerAvail: isInstallerExisted,
			config:                 config,
		}
		return newInstance, nil
	}
	return nil, nil
}

func (instance *CustomExec) execScript(script string) error {
	cmd := exec.Command(sh, script)
	if _, err := os.Stat(instance.config.GetInstallDir()); err == nil {
		cmd.Dir = instance.config.GetInstallDir()
	}
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}

func (instance CustomExec) StartPreInstall() error {
	// has pre install script?
	if _, err := os.Stat(instance.preInstall + "/" + constants.PREINSTALL_SH); err == nil {
		log.Println("Started Preinstall script")
		return instance.execScript(instance.preInstall)
	}
	return nil
}

func (instance CustomExec) StartPostInstall() error {
	// has pre install script?
	if _, err := os.Stat(instance.postInstall + "/" + constants.PREINSTALL_SH); err == nil {
		log.Println("Started Postinstall script")
		return instance.execScript(instance.preInstall)
	}
	return nil
}

func (instance CustomExec) RunCustomInstaller(srcFile string) error {
	cinstaller := getLocation() + "/" + customInstallerName
	log.Println("Start running custom installer ", cinstaller)
	cmd := exec.Command(cinstaller, srcFile)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}

func (instance CustomExec) Clean() error {
	loc := getLocation()
	// delete files
	if _, err := os.Stat(loc); err == nil {
		if err := os.RemoveAll(loc); err != nil {
			return err
		}
	}
	return nil
}
