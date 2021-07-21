package main

import (
	"ela/foundation/constants"
	"ela/foundation/event"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"log"
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
	case "help":
		println("Commands:")
		println("terminate/t", "-", "Terminate the current running system and its all apps.")
		println("status", "-", "Use to check the current status of system.")
	}
}

// use to terminate the system
// @timeout is the time it takes to terminate the system
func terminate(timeout int16) {
	log.Println("Terminating...")
	// step: check theres an existing connection with the system server.
	// if nothing then its terminated already
	con := connectToSystem()
	if con == nil {
		log.Println("System already terminated.")
		return
	}
	res, err := con.SendServiceRequest(
		constants.SYSTEM_SERVICE_ID,
		data.NewAction(constants.SYSTEM_TERMINATE, "", timeout))
	if err != nil {
		log.Println("Terminate ", err.Error())
		return
	}
	log.Println("Terminate", res)
}

// connect to system. return connector if success
func connectToSystem() protocol.ConnectorClient {
	con := event.CreateClientConnector()
	if err := con.Open(3); err != nil {
		return nil
	}
	return con
}

func getStatus() string {
	if connectToSystem() == nil {
		return "system was stopped"
	}
	return "system is running"
}
