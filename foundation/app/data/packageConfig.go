package data

import (
	"archive/zip"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const SYSTEM = "system"
const EXTERNAL = "external"
const NODE_JS = "nodejs"

type PackageConfig struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	PackageId        string   `json:"packageId"`      // company.package
	Build            int16    `json:"build"`          // major.minor.patch
	Version          string   `json:"version"`        // current version
	Permissions      []string `json:"permissions"`    // declared actions to be called inside the app
	ExportServices   bool     `json:"exportService"`  // true if the package contains services
	Activities       []string `json:"activities"`     // if app has activity. this contains definition of actions that will triggerr activity
	BroacastListener []string `json:"actionListener"` // defined actions which action listener will listen to
	InstallLocation  string   `json:"location"`       // which location the package will be installed
	Source           string   `json:"-"`              // the source location
	Nodejs           bool     `json:"nodejs"`         // true if this package includes node js
	//Restart          bool     `json:"restart"`        // if true system restart upon installation
	//Services         map[string]string `json:"services"`       // if app has a service. this contains definition of commands available to service
}

func DefaultPackage() *PackageConfig {
	return &PackageConfig{InstallLocation: EXTERNAL /*, Restart: false*/}
}

func (c *PackageConfig) LoadFromSrc(src string) error {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return &PackageConfigError{propertyError: "Error loading package." + err.Error()}

	}
	c.Source = src
	return c.LoadFromBytes(bytes)
}

func (c *PackageConfig) LoadFromReader(reader io.Reader) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return &PackageConfigError{propertyError: err.Error()}
	}
	err = c.LoadFromBytes(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *PackageConfig) LoadFromBytes(bytes []byte) error {
	return json.Unmarshal(bytes, &c)
}

// use to check if this package is valid
func (c *PackageConfig) IsValid() bool {
	return c.PackageId != ""
}

// returns true if this package is part of the system
func (c *PackageConfig) IsSystemPackage() bool {
	return c.InstallLocation == SYSTEM
}

// if current locatio is external. move it to system
func (c *PackageConfig) ChangeToSystemLocation() {
	c.InstallLocation = SYSTEM
}

// return true is this package contains services
func (c *PackageConfig) HasServices() bool {
	return c.ExportServices
}

// load package info from zip
func (c *PackageConfig) LoadFromZipPackage(source string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return errors.SystemNew("Load Package error. @"+source, err)
	}
	defer reader.Close()
	return c.LoadFromZipFiles(reader.File)
}

func (c *PackageConfig) LoadFromZipFiles(files []*zip.File) error {
	for _, file := range files {
		if file.Name != constants.APP_CONFIG_NAME {
			continue
		}
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		err = c.LoadFromReader(reader)
		if err != nil {
			return err
		}
		break
	}
	return nil
}

// return where this package should be installed
func (c *PackageConfig) GetInstallDir() string {
	if c.InstallLocation == SYSTEM || !path.HasExternal() {
		return path.GetSystemAppDir() + "/" + c.PackageId
	} else {
		return path.GetExternalAppDir() + "/" + c.PackageId
	}
}

func (c *PackageConfig) GetDataDir() string {
	if c.InstallLocation == SYSTEM || !path.HasExternal() {
		return path.GetSystemAppDirData(c.PackageId)
	} else {
		return path.GetExternalAppData(c.PackageId)
	}
}

func (c *PackageConfig) GetNodejsDir() string {
	return c.GetInstallDir() + "/" + NODE_JS
}

func (c *PackageConfig) GetMainExec() string {
	return path.GetAppMain(c.PackageId, c.InstallLocation == EXTERNAL)
}

func (c *PackageConfig) GetLibraryDir() string {
	return path.GetLibPath() + "/" + c.PackageId
}

func (c *PackageConfig) HasMainExec() bool {
	if _, err := os.Stat(c.GetMainExec()); err == nil {
		return true
	}
	return false
}

func (c *PackageConfig) ToString() string {
	json, _ := json.Marshal(c)
	return string(json)
}
