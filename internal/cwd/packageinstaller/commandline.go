package main

import (
	"ela/foundation/app/data"
	"ela/foundation/path"
	"ela/foundation/perm"
	"ela/internal/cwd/packageinstaller/utils"
	"log"
	"os"
)

/*
	commandline.go
	Commandline version of installer
*/

func startCommandline() {
	println("Elabox Installer Commandline")
	println("type help or -h for arguments.")

	// step: commandline help?
	if IsArgExist("help") || IsArgExist("-h") {
		println("usage:")
		println("command <path to package> -r(to restart the system)")
		return
	}

	// step: terminate the system first
	if err := utils.TerminateSystem(); err != nil {
		log.Println("Terminate system error", err)
	}
	// step: check if this is parent or not base on lock file
	isParentInstaller := true
	lockFile := path.GetCacheDir() + "/installer.lock"
	if _, err := os.Stat(lockFile); err == nil {
		isParentInstaller = false
	} else {
		file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, perm.PUBLIC)
		if err != nil {
			log.Println("Coulndt create @"+lockFile, err)
			return
		}
		defer os.Remove(lockFile)
		defer file.Close()
	}

	args := os.Args
	// step: check if valid path
	targetPk := args[1]
	if _, err := os.Stat(targetPk); err != nil {
		log.Fatal("Unable to install package with invalid path. "+targetPk, err)
		return
	}
	// load the package info
	pkg := data.DefaultPackage()
	if err := pkg.LoadFromZipPackage(targetPk); err != nil {
		log.Fatal("Unable to load package info", err)
		return
	}
	if !pkg.IsValid() {
		log.Fatalln("Package is not valid")
		return
	}
	newInstall := installer{BackupEnabled: true, PackageInfo: pkg, RunCustomInstaller: isParentInstaller}
	// step: start install
	if err := newInstall.Decompress(targetPk); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			log.Println("Failed reverting installer.", err.Error())
		}
		log.Fatal(err.Error())
		return
	}
	// step: post install
	if err := newInstall.Postinstall(); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			log.Println("Failed reverting installer.", err.Error())
		}
		log.Fatal(err.Error())
		return
	}
	log.Println("Installed success.")
	// step: restart system
	if IsArgExist("-r") {
		if err := utils.RestartSystem(); err != nil {
			log.Fatal(err.Error())
			return
		}
	}
}

func IsArgExist(arg string) bool {
	args := os.Args[1:]
	for _, _arg := range args {
		if arg == _arg {
			return true
		}
	}
	return false
}
