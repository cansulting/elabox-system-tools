package main

import (
	"archive/zip"
	"ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/internal/cwd/global"
	cwdg "ela/internal/cwd/global"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

/*
	package.go
	This struct handles the packaging.
	Config file needs to be loaded first before it can be used.
*/
type Package struct {
	Cwd             string   `json:"cwd"`             // current working directory
	PackageConfig   string   `json:"config"`          // the package config file
	BinDir          string   `json:"binDir"`          // bin directory. if bin file property was provided this will be skip
	Bin             string   `json:"bin"`             // bin file
	Lib             string   `json:"lib"`             // shared library for this package.
	Packages        []string `json:"packages"`        // list of packages to be included
	Www             string   `json:"www"`             // www front end to be included in source
	Nodejs          string   `json:"nodejs"`          // add node js directory if the package contain node js app
	PostInstall     string   `json:"postinstall"`     //
	PreInstall      string   `json:"preinstall"`      //
	CustomInstaller string   `json:"customInstaller"` // custom installer
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
	/*
		if c.Bin == "" && c.BinDir == "" && c.Www == "" && c.Nodejs == "" {
			return nil, errors.SystemNew("Package.loadPackageConfig() failed, bin/binDir/www/nodejs shouldnt be empty", nil)
		}*/
	if c.PackageConfig == "" {
		return nil, errors.SystemNew("Package.loadPackageConfig() failed, config shouldnt be empty", nil)
	}
	pkconfig := data.DefaultPackage()
	if err := pkconfig.LoadFromSrc(c.PackageConfig); err != nil {
		return nil, errors.SystemNew("Package.loadPackageConfig() failed", err)
	}
	if !pkconfig.IsValid() {
		return nil, errors.SystemNew("Package.loadPackageConfig() Is not valid", nil)
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
		return errors.SystemNew("Package.Compile() create zip file failed.", err)
	}
	zipwriter := zip.NewWriter(file)
	// add config
	if err := addFile(constants.APP_CONFIG_NAME, c.PackageConfig, zipwriter); err != nil {
		log.Println("Compile() adding package info")
		return errors.SystemNew("Package.Compile() add config file failed.", err)
	}
	// add binaries
	if c.Bin != "" {
		if err := addFile("bin/"+path.MAIN_EXEC_NAME, c.Bin, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding binary.", err)
		}
	} else if c.BinDir != "" {
		if err := addDir("bin", c.BinDir, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding binary dir.", err)
		}
	}
	// add shared libraries
	if c.Lib != "" {
		if err := addDir("lib", c.Lib, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding lib files @ "+c.Lib, err)
		}
	}
	// add packages
	if c.Packages != nil {
		for _, p := range c.Packages {
			log.Println("Compile() adding subpackage " + p)
			pkconfig := data.DefaultPackage()
			if err := pkconfig.LoadFromZipPackage(p); err != nil {
				return errors.SystemNew("Package.Compile() failed loading package "+p, err)
			}
			if err := addFile("packages/"+pkconfig.PackageId + "." + global.PACKAGE_EXT, p, zipwriter); err != nil {
				return errors.SystemNew("Package.Compile() failed adding package "+p, err)
			}
		}
	}
	// add www app
	if c.Www != "" {
		log.Println("Compile() adding wwww")
		if err := addDir(global.PKEY_WWW, c.Www, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding node js direcory @ "+c.Www, err)
		}
	}
	// add node js app
	if c.Nodejs != "" {
		log.Println("Compile() adding nodejs")
		if err := addDir("nodejs", c.Nodejs, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding node js direcory @ "+c.Nodejs, err)
		}
	}
	// add custom installer
	if c.CustomInstaller != "" {
		log.Println("Compile() adding custom installer")
		if err := addFile(
			cwdg.PACKAGEKEY_CUSTOM_INSTALLER+filepath.Ext(c.CustomInstaller),
			c.CustomInstaller, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding custom installer "+c.CustomInstaller, err)
		}
	}
	// scripts
	if c.PreInstall != "" {
		if err := addFile("scripts/"+cwdg.PREINSTALL_SH, c.PreInstall, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding script "+c.PreInstall, err)
		}
	}
	if c.PostInstall != "" {
		if err := addFile("scripts/"+cwdg.POSTINSTALL_SH, c.PostInstall, zipwriter); err != nil {
			return errors.SystemNew("Package.Compile() failed adding script "+c.PreInstall, err)
		}
	}
	// close
	if err := zipwriter.Close(); err != nil {
		return errors.SystemNew("Package.Compile() close failed.", err)
	}
	log.Println("Package success", pkconfig.PackageId)
	return nil
}

// add file
// @name: header for zip file
// @src: location of the file so it can be read
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
