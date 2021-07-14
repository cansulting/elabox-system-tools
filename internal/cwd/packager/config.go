package main

import (
	"archive/zip"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Cwd             string   `json:"cwd"`             // current working directory
	PackageConfig   string   `json:"config"`          // the package config file
	BinDir          string   `json:"binDir"`          // bin directory. if bin file property was provided this will be skip
	Bin             string   `json:"bin"`             // bin file
	Packages        []string `json:"packages"`        // list of packages to be included
	CustomInstaller string   `json:"customInstaller"` // custom installer
	Www             string   `json:"www"`             // www front end to be included in source
	Nodejs          string   `json:"nodejs"`          // add node js directory if the package contain node js app
}

// load packager config file from path
func (c *Config) LoadFrom(src string) error {
	log.Println("Packaging @", src)
	bytes, err := os.ReadFile(src)
	if err != nil {
		return errors.SystemNew("Config.LoadFrom() failed "+src, err)
	}
	if err := json.Unmarshal(bytes, c); err != nil {
		return errors.SystemNew("Config.LoadFrom() failed "+src, err)
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
func (c *Config) loadPackageConfig() (*data.PackageConfig, error) {
	if c.Bin == "" && c.BinDir == "" && c.Www == "" && c.Nodejs == "" {
		return nil, errors.SystemNew("Config.loadPackageConfig() failed, bin/binDir/www/nodejs shouldnt be empty", nil)
	}
	if c.PackageConfig == "" {
		return nil, errors.SystemNew("Config.loadPackageConfig() failed, config shouldnt be empty", nil)
	}
	pkconfig := data.DefaultPackage()
	if err := pkconfig.LoadFromSrc(c.PackageConfig); err != nil {
		return nil, errors.SystemNew("Config.loadPackageConfig() failed", err)
	}
	if !pkconfig.IsValid() {
		return nil, errors.SystemNew("Config.loadPackageConfig() Is not valid", nil)
	}
	return pkconfig, nil
}

// use to start compiling
// @destDir path to where file will be save
func (c *Config) Compile(destdir string) error {
	pkconfig, err := c.loadPackageConfig()
	if err != nil {
		return err
	}
	outputp := destdir + "/" + pkconfig.PackageId + "." + constants.PACKAGE_EXT
	log.Println("Output path", outputp)
	// create zip file
	file, err := os.Create(outputp)
	if err != nil {
		return errors.SystemNew("Config.Compile() create zip file failed.", err)
	}
	zipwriter := zip.NewWriter(file)
	// add config
	if err := addFile(constants.APP_CONFIG_NAME, c.PackageConfig, zipwriter); err != nil {
		return errors.SystemNew("Config.Compile() add config file failed.", err)
	}
	// add binaries
	if c.Bin != "" {
		if err := addFile("bin/"+path.MAIN_EXEC_NAME, c.Bin, zipwriter); err != nil {
			return errors.SystemNew("Config.Compile() failed adding binary.", err)
		}
	}
	// add packages
	if c.Packages != nil {
		for _, p := range c.Packages {
			pkconfig := data.DefaultPackage()
			if err := pkconfig.LoadFromZipPackage(p); err != nil {
				return errors.SystemNew("Config.Compile() failed loading package "+p, err)
			}
			if err := addFile("packages/"+pkconfig.PackageId, p, zipwriter); err != nil {
				return errors.SystemNew("Config.Compile() failed adding package "+p, err)
			}
		}
	}
	// add www app
	if c.Www != "" {
		if err := addDir("www", c.Www, zipwriter); err != nil {
			return errors.SystemNew("Config.Compile() failed adding node js direcory @ "+c.Nodejs, err)
		}
	}
	// add node js app
	if c.Nodejs != "" {
		if err := addDir("nodejs", c.Nodejs, zipwriter); err != nil {
			return errors.SystemNew("Config.Compile() failed adding node js direcory @ "+c.Nodejs, err)
		}
	}
	// add custom installer
	if c.CustomInstaller != "" {
		if err := addFile("installer/installer"+
			filepath.Ext(c.CustomInstaller), c.CustomInstaller, zipwriter); err != nil {
			return errors.SystemNew("Config.Compile() failed adding custom installer "+c.CustomInstaller, err)
		}
	}
	// close
	if err := zipwriter.Close(); err != nil {
		return errors.SystemNew("Config.Compile() close failed.", err)
	}
	return nil
}

func addFile(name string, src string, w *zip.Writer) error {
	f, err := w.Create(name)
	if err != nil {
		return err
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
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := addFile(
			newDirName+"/"+file.Name(),
			srcDir+"/"+file.Name(),
			w); err != nil {
			return err
		}
	}
	return nil
}
