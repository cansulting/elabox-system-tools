package utils

import (
	"archive/zip"
	"ela/foundation/app/data"
	"ela/foundation/errors"
	"ela/foundation/perm"
	cwdg "ela/internal/cwd/global"
	"ela/internal/cwd/packageinstaller/global"
	"log"
	"os"
	"os/exec"
)

/*
	CustomExec.go
	this struct handles execution of scripts/sh and custom installer attach to it
*/

const sh = "/bin/sh"

type CustomExec struct {
	config      *data.PackageConfig
	preInstall  string // temp path of preinstall script
	postInstall string // temp path of post install script
}

func GetSHConfig(config *data.PackageConfig, files []*zip.File) (*CustomExec, error) {
	cinstallerDir := global.GetTempPath()

	// export installer and scripts
	filters := []Filter{
		{Keyword: cwdg.PKEY_POST_INSTALL_SH, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
		{Keyword: cwdg.PKEY_PRE_INSTALL_SH, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
		{Keyword: cwdg.PACKAGEKEY_CUSTOM_INSTALLER, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
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
		newInstance := &CustomExec{
			preInstall:  cinstallerDir + "/" + cwdg.PKEY_PRE_INSTALL_SH,
			postInstall: cinstallerDir + "/" + cwdg.PKEY_POST_INSTALL_SH,
			config:      config,
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
	cmd.Stdout = instance
	cmd.Stderr = instance
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// start running pre install script
func (instance CustomExec) StartPreInstall() error {
	// has pre install script?
	if _, err := os.Stat(instance.preInstall); err == nil {
		log.Println("-------------Execute pre install script-------------")
		return instance.execScript(instance.preInstall)
	}
	return nil
}

// start running post install script
func (instance CustomExec) StartPostInstall() error {
	// has pre install script?
	if _, err := os.Stat(instance.postInstall); err == nil {
		log.Println("-------------Execute post install script-------------")
		return instance.execScript(instance.postInstall)
	}
	return nil
}

// run custom installer
func (instance CustomExec) RunCustomInstaller(srcFile string) error {
	cinstaller := global.GetCustomInstallerTempPath()
	log.Println("Start running custom installer ", cinstaller)
	cmd := exec.Command(cinstaller, srcFile)
	cmd.Stdout = instance
	cmd.Stderr = instance
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// called when custom installer has write log
func (instance CustomExec) Write(bytes []byte) (int, error) {
	print(string(bytes))
	return len(bytes), nil
}

func (instance CustomExec) Clean() error {
	loc := global.GetTempPath()
	// delete files
	if _, err := os.Stat(loc); err == nil {
		if err := os.RemoveAll(loc); err != nil {
			return err
		}
	}
	return nil
}
