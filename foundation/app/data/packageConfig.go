// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file handles functions and producedure related to Package configuration

package data

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

const SYSTEM = "system"      // identifies the package is installed on system location
const EXTERNAL = "external"  // identifies the package is installed on external location
const NODE_JS_DIR = "nodejs" // sub directory of binary dir, this is where node js scripts reside

// This structure represents package  json file along with the binary.
// this contains information about the application behaviour, permission and services.
type PackageConfig struct {
	Name        string `json:"name"`        // package name
	Description string `json:"description"` // description of package
	PackageId   string `json:"packageId"`   // identifies the package/application. this should be unique. format = company.package
	Build       int16  `json:"build"`       // this should be incremental starting from 1
	Version     string `json:"version"`     // major.minor.patch
	Program     string `json:"program"`     // the main program file to execute
	// request permission for specific action/feature
	// if the specific action was called and was not defined. the process will be void
	Permissions      []string               `json:"permissions"`
	ExportServices   bool                   `json:"exportService"`   // true if the package contains services
	Activities       []string               `json:"activities"`      // if app has activity. this contains definition of actions that will triggerr activity
	BroacastListener []string               `json:"actionListener"`  // defined actions which action listener will listen to
	InstallLocation  string                 `json:"location"`        // which location the package will be installed
	Source           string                 `json:"-"`               // the source location
	Nodejs           bool                   `json:"nodejs"`          // true if this package includes node js
	PackagerVersion  string                 `json:"packagerVersion"` // version of packager of this package
	Ext              map[string]interface{} `json:"ext"`             // extra values
	//Services         map[string]string `json:"services"`       // if app has a service. this contains definition of commands available to service
}

// default values for package
func DefaultPackage() *PackageConfig {
	return &PackageConfig{InstallLocation: EXTERNAL /*, Restart: false*/}
}

// local package data given the source location
func (c *PackageConfig) LoadFromSrc(src string) error {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return errors.SystemNew("Error loading package.", err)

	}
	c.Source = src
	return c.LoadFromBytes(bytes)
}

// load package based from reader info
func (c *PackageConfig) LoadFromReader(reader io.Reader) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.SystemNew("Failed loading package from reader", err)
	}
	err = c.LoadFromBytes(bytes)
	if err != nil {
		return err
	}
	return nil
}

// load package base from bytes
func (c *PackageConfig) LoadFromBytes(bytes []byte) error {
	return json.Unmarshal(bytes, &c)
}

// use to check if this package is valid
func (c *PackageConfig) GetIssue() (string, string) {
	if c.PackageId == "" {
		return "packageId", "Input a valid packageId.  eg <company>.<app name>"
	}
	if c.Name == "" {
		return "name", "Provide a proper name for package."
	}
	// if !c.Nodejs && c.Program == "" {
	// 	return "program", "Provide a valid file name for main program entry."
	// }
	if c.Build < 0 {
		return "build", "Provide a valid build number. Value should be greater to 0"
	}
	return "", ""
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

// use to check if contains activity that has action of
func (c *PackageConfig) HasActivity(actionId string) bool {
	if c.Activities == nil || len(c.Activities) == 0 {
		return false
	}
	for _, act := range c.Activities {
		if act == actionId {
			return true
		}
	}
	return false
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

// lookup package from base from zip files
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

// get data directory of package
func (c *PackageConfig) GetDataDir() string {
	if c.InstallLocation == SYSTEM || !path.HasExternal() {
		return path.GetSystemAppDirData(c.PackageId)
	} else {
		return path.GetExternalAppData(c.PackageId)
	}
}

// get package's node js dir
func (c *PackageConfig) GetNodejsDir() string {
	return c.GetInstallDir() + "/" + NODE_JS_DIR
}

// get package's main binary
// returns binary location
func (c *PackageConfig) GetMainProgram() string {
	return path.GetAppInstallLocation(c.PackageId, c.InstallLocation == EXTERNAL) + "/" + c.Program
}

// get package's library directory
func (c *PackageConfig) GetLibraryDir() string {
	return path.GetLibPath() + "/" + c.PackageId
}

// return true if has main binary
func (c *PackageConfig) HasMainProgram() bool {
	if c.Program == "" {
		return false
	}
	if _, err := os.Stat(c.GetMainProgram()); err == nil {
		return true
	}
	return false
}

// return json information
func (c *PackageConfig) ToJson() string {
	json, _ := json.MarshalIndent(c, "", "\t")
	return string(json)
}
