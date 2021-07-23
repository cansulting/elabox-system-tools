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
	// true if restarts system
	restartSystem := IsArgExist("-r")
	// step: check if this is parent or not base on lock file
	isParentInstaller := true
	// if restartSystem is true, automatically run as parent
	if !restartSystem {
		lockFile := getLockFile()
		if _, err := os.Stat(lockFile); err == nil {
			log.Println("Lock file exist, this is child installer.")
			isParentInstaller = false
		} else {
			file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, perm.PUBLIC)
			if err != nil {
				log.Println("Coulndt create @"+lockFile, err)
				return
			}
			file.Close()
			defer removeLockFile()
		}
	}
	// step: terminate the system first
	if isParentInstaller {
		if err := utils.TerminateSystem(); err != nil {
			log.Println("Terminate system error", err)
		}
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
		removeLockFile()
		log.Fatal(err.Error())
		return
	}
	removeLockFile()
	// step: post install
	if err := newInstall.Finalize(); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			log.Println("Failed reverting installer.", err.Error())
		}
		log.Fatal(err.Error())
		return
	}
	log.Println("Installed success.")
	// step: restart system
	if restartSystem {
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

func getLockFile() string {
	return path.GetCacheDir() + "/installer.lock"
}

func removeLockFile() {
	if _, err := os.Stat(getLockFile()); err == nil {
		os.Remove(getLockFile())
	}
}
