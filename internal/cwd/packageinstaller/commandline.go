package main

import (
	"ela/foundation/app/data"
	"ela/internal/cwd/packageinstaller/utils"
	"log"
	"os"
)

/*
	commandline.go
	Commandline version of installer
*/

func startCommandline() {
	log.Println("Elabox Installer Commandline")
	args := os.Args
	// step: check if valid path
	packagePath := args[1]
	if _, err := os.Stat(packagePath); err != nil {
		log.Fatal("Unable to install package with invalid path. "+packagePath, packagePath)
		return
	}
	// load the package info
	pkg := data.DefaultPackage()
	if err := pkg.LoadFromZipPackage(packagePath); err != nil {
		log.Fatal("Unable to load package info", err)
		return
	}
	if !pkg.IsValid() {
		log.Fatalln("Package is not valid")
		return
	}
	newInstall := installer{BackupEnabled: true, PackageInfo: pkg}
	// step: start install
	if err := newInstall.Decompress(packagePath); err != nil {
		log.Fatal(err.Error())
		return
	}
	// step: register package
	if err := newInstall.RegisterPackage(); err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println("Installed success.")
	// step: restart system
	if err := utils.RestartSystem(); err != nil {
		log.Fatal(err.Error())
		return
	}
}
