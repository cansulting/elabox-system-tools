package main

import (
	"ela/foundation/constants"
	"ela/foundation/event"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/foundation/system"
	"ela/internal/cwd/system/config"
	"os"
)

// process commandline
func processCmdline() {
	args := os.Args
	cmd := args[1]
	switch cmd {
	case "terminate", "-t":
		terminate(0)
	case "status":
		println(getStatus())
	case "env":
		print(config.GetEnv(args[2]))
	case "help":
		println("Commands:")
		println("terminate/t", "-", "Terminate the current running system and its all apps.")
		println("status", "-", "Use to check the current status of system.")
		println("env", "-", "Use to set or get environment variable")
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
	res, err := con.SendServiceRequest(
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
	return system.GetStatus()
}
