package main

import (
	"archive/zip"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/internal/cwd/global"
)

const VERSION = "0.1.0"

/*
	package.go
	This struct handles the packaging.
	Config file needs to be loaded first before it can be used.
*/
type Package struct {
	Cwd               string   `json:"cwd"`               // current working directory
	PackageConfig     string   `json:"config"`            // the package config file
	BinDir            string   `json:"binDir"`            // bin directory. if bin file property was provided this will be skip
	Bin               string   `json:"bin"`               // bin file
	Lib               string   `json:"lib"`               // shared library for this package.
	Packages          []string `json:"packages"`          // list of packages to be included
	Www               string   `json:"www"`               // www front end to be included in source
	Nodejs            string   `json:"nodejs"`            // add node js directory if the package contain node js app
	PostInstall       string   `json:"postinstall"`       // script that will be executed after everything is installed
	PreInstall        string   `json:"preinstall"`        // script that will be executed upon initialization
	UnInstall         string   `json:"uninstall"`         // script that will be called whenever uninstalling a package
	CustomInstaller   string   `json:"customInstaller"`   // custom installer
	ContinueOnMissing bool     `json:"continueOnMissing"` // true if continue to package if theres a missing subpackage
}

func NewPackage() *Package {
	return &Package{
		ContinueOnMissing: false,
	}
}

// load packager config file from path
func (c *Package) LoadConfig(src string) error {
	log.Println("Reading packager config @", src)
	bytes, err := os.ReadFile(src)
	if err != nil {
		return errors.SystemNew("Package.LoadFrom() failed "+src, err)
	}
	if err := json.Unmarshal(bytes, c); err != nil {
		return errors.SystemNew("Package.LoadFrom() failed "+src, err)
	}
	// change working directory
	cwd := c.Cwd
	if cwd == "" {
		cwd = filepath.Dir(filepath.Clean(src))
	}
	os.Chdir(cwd)
	log.Println("Changed working directory to", cwd)
	return nil
}

// use to loadPackageConfig parameters before packaging
func (c *Package) loadPackageConfig() (*data.PackageConfig, error) {
	if c.PackageConfig == "" {
		return nil, errors.SystemNew("loadPackageConfig() failed, config shouldnt be empty", nil)
	}
	pkconfig := data.DefaultPackage()
	if err := pkconfig.LoadFromSrc(c.PackageConfig); err != nil {
		return nil, errors.SystemNew("loadPackageConfig() failed", err)
	}
	// check if theres issue with package config
	issueProperty, issueMsg := pkconfig.GetIssue()
	if issueProperty != "" {
		return nil, errors.SystemNew(
			"loadPackageConfig(): Package "+pkconfig.Source+" has issues. Property "+issueProperty+" - "+issueMsg,
			nil)
	}
	return pkconfig, nil
}

// use to start compiling
// @destDir path to where file will be save
func (c *Package) Compile(destdir string) error {
	pkconfig, err := c.loadPackageConfig()
	if err != nil {
		return err
	}
	outputp := destdir + "/" + pkconfig.PackageId + "." + constants.PACKAGE_EXT
	log.Println("Output path", outputp)
	// create zip file
	file, err := os.Create(outputp)
	if err != nil {
		return errors.SystemNew("Compile() create zip file failed.", err)
	}
	zipwriter := zip.NewWriter(file)
	// add config
	if err := addFile(constants.APP_CONFIG_NAME, c.PackageConfig, zipwriter); err != nil {
		log.Println("Compile() adding package info")
		return errors.SystemNew("Compile() add config file failed.", err)
	}
	// add binaries
	if c.Bin != "" {
		if err := addFile("bin/"+pkconfig.Program, c.Bin, zipwriter); err != nil {
			return errors.SystemNew("Failed adding binary.", err)
		}
	} else if c.BinDir != "" {
		// check if the entry program exist in the directory
		if pkconfig.Program != "" {
			if _, err := os.Stat(c.BinDir + "/" + pkconfig.Program); err != nil {
				return errors.SystemNew("Main program not found @ "+c.BinDir+"/"+pkconfig.Program+".", nil)
			}
		}
		if err := addDir("bin", c.BinDir, zipwriter); err != nil {
			return errors.SystemNew("Failed adding binary dir.", err)
		}
	} //else {
	// 	return errors.SystemNew("Failed. No binaries were provided.", nil)
	// }
	// add shared libraries
	if c.Lib != "" {
		if err := addDir("lib", c.Lib, zipwriter); err != nil {
			return errors.SystemNew("failed adding lib files @ "+c.Lib, err)
		}
	}
	// add sub packages
	if c.Packages != nil {
		for _, p := range c.Packages {
			log.Println("Compile() adding subpackage " + p)
			pkconfig := data.DefaultPackage()
			if err := pkconfig.LoadFromZipPackage(p); err != nil {
				if !c.ContinueOnMissing {
					return errors.SystemNew("failed loading subpackage package "+p, err)
				} else {
					log.Println("Warning: subpackage cannot be found @ ", p, ". skipped.")
					continue
				}
			}
			if err := addFile("packages/"+pkconfig.PackageId+"."+global.PACKAGE_EXT, p, zipwriter); err != nil {
				return errors.SystemNew("failed adding package "+p, err)
			}
		}
	}
	// add www app
	if c.Www != "" {
		log.Println("Compile() adding wwww")
		if err := addDir(global.PKEY_WWW, c.Www, zipwriter); err != nil {
			return errors.SystemNew("failed adding www direcory @ "+c.Www, err)
		}
	}
	// add node js app
	if c.Nodejs != "" {
		log.Println("Compile() adding nodejs")
		if err := addDir("nodejs", c.Nodejs, zipwriter); err != nil {
			return errors.SystemNew("failed adding node js direcory @ "+c.Nodejs, err)
		}
	}
	// add custom installer
	if c.CustomInstaller != "" {
		log.Println("adding custom installer")
		if err := addFile(
			global.PACKAGEKEY_CUSTOM_INSTALLER+filepath.Ext(c.CustomInstaller),
			c.CustomInstaller, zipwriter); err != nil {
			return errors.SystemNew("failed adding custom installer "+c.CustomInstaller, err)
		}
	}
	// scripts
	if c.PreInstall != "" {
		if err := addFile("scripts/"+global.PREINSTALL_SH, c.PreInstall, zipwriter); err != nil {
			return errors.SystemNew("failed adding script "+c.PreInstall, err)
		}
	}
	if c.PostInstall != "" {
		if err := addFile("scripts/"+global.POSTINSTALL_SH, c.PostInstall, zipwriter); err != nil {
			return errors.SystemNew("failed adding script "+c.PreInstall, err)
		}
	}
	if c.UnInstall != "" {
		if err := addFile("bin/"+global.UNINSTALL_SH, c.UnInstall, zipwriter); err != nil {
			return errors.SystemNew("failed adding script "+c.UnInstall, err)
		}
	}
	// close
	if err := zipwriter.Close(); err != nil {
		return errors.SystemNew("Compile() close failed.", err)
	}
	// update package info version
	pkconfig.PackagerVersion = VERSION
	if err := os.WriteFile(pkconfig.Source, []byte(pkconfig.ToJson()), perm.PUBLIC); err != nil {
		return errors.SystemNew("Updating package config "+pkconfig.Source+"Failed", err)
	}
	log.Println("Package success", pkconfig.PackageId)
	return nil
}

// add file
// @name: new path/name to zip
// @src: location of the file so it can be read
func addFile(name string, src string, w *zip.Writer) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	if _, err := os.Stat(src); err != nil {
		log.Println(src + " skipped. File not found")
		return nil
	}
	bytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	_, err2 := f.Write(bytes)
	if err != nil {
		return err2
	}
	return nil
}

func addDir(newDirName string, srcDir string, w *zip.Writer) error {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Type() == fs.ModeDir {
			if err := addDir(
				newDirName+"/"+file.Name(),
				srcDir+"/"+file.Name(),
				w); err != nil {
				return err
			}
		} else {
			if err := addFile(
				newDirName+"/"+file.Name(),
				srcDir+"/"+file.Name(),
				w); err != nil {
				return err
			}
		}
	}
	return nil
}
