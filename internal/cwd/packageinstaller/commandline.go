package main

import (
	"ela/foundation/constants"
	"ela/foundation/event"
	"ela/foundation/path"
	"ela/foundation/perm"
	"ela/foundation/system"
	"ela/internal/cwd/global/server"
	"ela/internal/cwd/packageinstaller/landing"
	"ela/internal/cwd/packageinstaller/pkg"
	"ela/internal/cwd/packageinstaller/utils"
	"ela/internal/cwd/system/global"
	"log"
	"os"
	"time"
)

/*
	commandline.go
	Commandline version of installer
*/

type loghandler struct {
	logfile *os.File
}

func (i loghandler) init() {
	// create log file
	if IsArgExist("-l") {
		logp := path.GetCacheDir() + "/installer.log"
		log.Println("Check log @", logp)
		var err error
		i.logfile, err = os.OpenFile(logp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm.PUBLIC)
		if err == nil {
			log.SetOutput(i)
		}
	}
}
func (i loghandler) close() {
	i.logfile.Close()
}
func (i loghandler) Write(data []byte) (int, error) {
	print(string(data))
	if global.Connector != nil {
		global.Connector.Broadcast(constants.SYSTEM_SERVICE_ID, "log", string(data))
	}
	if i.logfile != nil {
		i.logfile.Write(data)
	}
	return len(data), nil
}

func startCommandline() {
	println("Elabox Installer Commandline")
	println("type help or -h for arguments.")

	// step: commandline help?
	if IsArgExist("help") || IsArgExist("-h") {
		println("usage:")
		println("Example command <path to package> -r -s")
		println("-r - to restart the system")
		println("-s - this is system update")
		println("-l - create log file")
		println("-i - ignore custom installer")
		return
	}

	// true if restarts system
	restartSystem := IsArgExist("-r") || IsArgExist("-s")
	// step: terminate the system first
	/*
		if restartSystem {
			if err := utils.TerminateSystem(); err != nil {
				log.Println("Terminate system error", err)
			}
		}*/
	args := os.Args
	targetPk := args[1]
	// step: load package
	content, err := pkg.LoadFromSource(targetPk)
	if err != nil {
		log.Fatal("Failed running commandline", err)
		return
	}
	// create logger?
	if IsArgExist("-l") || IsArgExist("-s") {
		logger := loghandler{}
		logger.init()
		defer logger.close()
	}
	// step: will initialize server?
	if IsArgExist("-s") {
		startServer(content)
	}
	// use custom installer or not?
	if IsArgExist("-i") || !content.HasCustomInstaller() {
		normalInstall(content)
	} else {
		if err := content.RunCustomInstaller(targetPk, true, "-i"); err != nil {
			log.Fatal("Failed running custom installer", err)
			return
		}
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

func normalInstall(content *pkg.Data) {
	// step: wait and make sure system was terminated. for system updates
	time.Sleep(time.Second)
	newInstall := NewInstaller(content, true)
	// step: start install
	if err := newInstall.Start(); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			log.Println("Failed reverting installer.", err.Error())
		}
		log.Fatal(err.Error())
		return
	}
	// step: post install
	if err := newInstall.Finalize(); err != nil {
		// failed? revert changes
		if err := newInstall.RevertChanges(); err != nil {
			log.Println("Failed reverting installer.", err.Error())
		}
		log.Fatal(err.Error())
		return
	}
}

// start installer server
func startServer(content *pkg.Data) {
	// serve socket io
	conn := event.CreateServerConnector()
	if err := conn.Open(); err != nil {
		log.Fatal("Failed to initialize intaller server.", err.Error())
		return
	}
	global.Connector = conn
	server.InitSystemService(conn, nil)
	// retrieve landing page first
	landingDir, err := content.ExtractLandingPage()
	// is there a landing page?
	if err == nil {
		if err := landing.Initialize(landingDir); err != nil {
			log.Fatal("Unable to initialize server.", err)
			return
		}
	} else {
		log.Println("Failed extracting landing page", err, ". Skipping www listener")
	}
	conn.SetStatus(system.UPDATING, nil)
	// step: if theres a landing page. wait for user to connect to landing page before continuing
	if landingDir != "" {
		landing.WaitForConnection()
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
