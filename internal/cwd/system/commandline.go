// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file provides commandline for system
// This can be use via
// sudo ebox -t        					-- to terminate
// sudo ebox -status   					-- to view system status
// sudo ebox -env <elabox environment>	-- to view specific elabox environment

package main

import (
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/env"
	"github.com/cansulting/elabox-system-tools/foundation/event"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

// process commandline
func processCmdline() {
	args := os.Args
	cmd := args[1]
	switch cmd {
	case "terminate", "-t":
		terminate(0)
	// use to status
	case "status", "-s":
		println(getStatus())
	// use to setting env
	case "env", "-e":
		largs := len(args)
		if largs == 3 {
			println(env.GetEnv(args[2]))
		} else if largs == 4 {
			env.SetEnv(args[2], args[3])
		}
	// use for version
	case "version", "-v":
		pkg, _ := app.RetrievePackage(constants.SYSTEM_SERVICE_ID)
		println(system.BuildMode, pkg)
	default:
		println("Commands:")
		println("terminate/-t", "-", "Terminate the current running system and its all apps.")
		println("status", "-s", "Use to check the current status of system.")
		println("env", "-e", "Use to set or get environment variable. eg. to get = env <name>, to set = env <name> <value>")
		println("version/-v", "-", "Check the current version.")
	}
}

// use to terminate the system
// @timeout is the time it takes to terminate the system
func terminate(timeout int16) {
	println("Terminating...")
	// step: check theres an existing connection with the system server.
	// if nothing then its terminated already
	con := connectToSystem()
	if con == nil {
		println("System already terminated.")
		return
	}
	res, err := con.SendSystemRequest(
		constants.SYSTEM_SERVICE_ID,
		data.NewAction(constants.SYSTEM_TERMINATE, "", timeout))
	if err != nil {
		println("Terminate ", err.Error())
		return
	}
	println("Terminate", res)
}

// connect to system. return connector if success
func connectToSystem() protocol.ConnectorClient {
	con := event.CreateClientConnector()
	if err := con.Open(1); err != nil {
		return nil
	}
	return con
}

func getStatus() string {
	return string(system.GetStatus())
}
