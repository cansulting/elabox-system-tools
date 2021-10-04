package pkg

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/internal/cwd/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/landing"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/utils"
)

const sh = "/bin/bash"

type Data struct {
	Config         *data.PackageConfig
	Files          []*zip.File
	customInstaler *zip.File
	preInstall     string // temp path of preinstall script
	postInstall    string // temp path of post install script
	zipInstance    *zip.ReadCloser
}

// load package from source
func LoadFromSource(src string) (*Data, error) {
	// step: validate action
	if src == "" {
		return nil, errors.SystemNew("Failed to load package. @"+src, nil)
	}
	src = filepath.Clean(src)
	constants.Logger.Debug().Msg("Loading package @" + src)
	// step: read package
	z, err := zip.OpenReader(src)
	if err != nil {
		return nil, errors.SystemNew("pkg.Load() failed to locate "+src, err)
	}
	res, err := LoadFromZipFiles(z.File)
	if res != nil {
		res.zipInstance = z
	}
	return res, err
}

func LoadFromZipFiles(files []*zip.File) (*Data, error) {
	res := &Data{}
	res.Files = files
	// step: init pkg config
	res.Config = data.DefaultPackage()
	if err := res.Config.LoadFromZipFiles(res.Files); err != nil {
		return nil, errors.SystemNew("Pkg.Load() failed loading config", err)
	}
	if !res.Config.IsValid() {
		return nil, errors.SystemNew("LoadFromZipFiles() Package is not valid", nil)
	}
	return res, nil
}

// is this package has custom installer?
func (instance *Data) HasCustomInstaller() bool {
	for _, file := range instance.Files {
		if file.Name == global.PACKAGEKEY_CUSTOM_INSTALLER {
			instance.customInstaler = file
			return true
		}
	}
	return false
}

// extract and run custom installer
func (instance *Data) RunCustomInstaller(srcPkg string, wait bool, args ...string) error {
	// step: lookup for installer
	if instance.customInstaler == nil {
		if !instance.HasCustomInstaller() {
			return errors.SystemNew("RunCustomInstaller() no custom installer found", nil)
		}
	}
	// step: open zip file
	reader, err := instance.customInstaler.Open()
	if err != nil {
		return errors.SystemNew("RunCustomInstaller failed copying from zip", err)
	}
	// step: extract and copy
	installerPath := constants.GetTempPath() + "/" + global.PACKAGEKEY_CUSTOM_INSTALLER
	if _, err := os.Stat(installerPath); err == nil {
		os.Remove(installerPath)
	}
	if err := utils.CopyToTarget(installerPath, reader, perm.PUBLIC); err != nil {
		return errors.SystemNew("RunCustomInstaller() failed cloning to tempfile", err)
	}
	// step: run
	constants.Logger.Debug().Msg("Running custom installer " + installerPath)
	args2 := make([]string, 0, 10)
	args2 = append(args2, srcPkg)
	if len(args) > 0 {
		args2 = append(args2, args...)
	}
	cmd := exec.Command(installerPath, args2...)
	//cmd.Dir = filepath.Dir(installerPath)
	cmd.Stdout = instance
	cmd.Stderr = instance
	cmd.Stdin = os.Stdin
	if wait {
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		if err := cmd.Start(); err != nil {
			return err
		}
		cmd.Process.Release()
	}
	return nil
}

// called when custom installer has write log
func (instance *Data) Write(values []byte) (int, error) {
	msg := string(values)
	print(msg)
	decoder := json.NewDecoder(bytes.NewReader(values))
	var log map[string]interface{}
	if err := decoder.Decode(&log); err != nil {
		landing.BroadcastLog(msg)
	} else {
		landing.BroadcastLog(log["message"].(string))
	}

	return len(values), nil
}

func (instance *Data) ExtractScripts() error {
	cinstallerDir := constants.GetTempPath()

	// export installer and scripts
	filters := []utils.Filter{
		{Keyword: global.PKEY_POST_INSTALL_SH, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
		{Keyword: global.PKEY_PRE_INSTALL_SH, InstallTo: cinstallerDir, Perm: perm.PUBLIC},
	}
	found := false
	// initialize scripts and files and move to cache directory
	for _, file := range instance.Files {
		reader, err := file.Open()
		if err != nil {
			return errors.SystemNew("ExtractScripts() Failed opening file", err)
		}
		defer reader.Close()
		for _, cfilter := range filters {
			newPath, err, applied := cfilter.CanApply(file.Name, reader, file.CompressedSize64)
			if err != nil {
				return err
			}
			if applied {
				reader, err := file.Open()
				if err != nil {
					return err
				}
				if cfilter.Save(newPath, reader); err != nil {
					return errors.SystemNew("Failed saving "+file.Name, err)
				}
				reader.Close()
				found = true
			}
		}
	}
	if found {
		instance.preInstall = cinstallerDir + "/" + global.PKEY_PRE_INSTALL_SH
		instance.postInstall = cinstallerDir + "/" + global.PKEY_POST_INSTALL_SH
	}
	return nil
}

// use to extract files to target path
func (instance *Data) ExtractFiles(keyword string, targetPath string, limit uint) (uint, error) {
	var found uint = 0
	for _, file := range instance.Files {
		if strings.Contains(file.Name, keyword) {
			found++
			reader, err := file.Open()
			if err != nil {
				return 0, err
			}
			defer reader.Close()
			fname := targetPath + "/" + file.Name
			if err := utils.CopyToTarget(fname, reader, perm.PUBLIC); err != nil {
				return 0, err
			}
			if found >= limit {
				return found, nil
			}
		}
	}
	return found, nil
}

func (instance *Data) HasPreInstallScript() bool {
	return instance.preInstall != ""
}

func (instance *Data) HasPostInstallScript() bool {
	return instance.postInstall != ""
}

func (instance *Data) execScript(script string) error {
	cmd := exec.Command(sh, script)
	if _, err := os.Stat(instance.Config.GetInstallDir()); err == nil {
		cmd.Dir = instance.Config.GetInstallDir()
	}
	cmd.Stdout = instance
	cmd.Stderr = instance
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// start running pre install script
func (instance *Data) StartPreInstall() error {
	// has pre install script?
	if _, err := os.Stat(instance.preInstall); err == nil {
		constants.Logger.Debug().Msg("-------------Execute pre install script-------------")
		return instance.execScript(instance.preInstall)
	}
	return nil
}

// start running post install script
func (instance *Data) StartPostInstall() error {
	// has pre install script?
	if _, err := os.Stat(instance.postInstall); err == nil {
		constants.Logger.Debug().Msg("-------------Execute post install script-------------")
		return instance.execScript(instance.postInstall)
	}
	return nil
}

func (instance *Data) Clean() error {
	loc := constants.GetTempPath()
	// delete files
	if _, err := os.Stat(loc); err == nil {
		if err := os.RemoveAll(loc); err != nil {
			return err
		}
	}
	return nil
}

// use to extract landing page
// returns path to the landing page dir hence return error if there are issues
func (instance *Data) ExtractLandingPage() (string, error) {
	targetp := constants.GetTempPath() + "/" + global.PKEY_WWW
	constants.Logger.Debug().Msg("Extracting landing @" + targetp)
	found, err := instance.ExtractFiles(global.PKEY_WWW, constants.GetTempPath(), 1000)
	if err != nil {
		return "", errors.SystemNew("Failed extracting landing page.", err)
	}
	if found <= 0 {
		return "", errors.SystemNew("Nothing was found on package with keyword "+global.PKEY_WWW, nil)
	}
	return targetp, nil
}

func (instance *Data) Close() {
	instance.zipInstance.Close()
}
