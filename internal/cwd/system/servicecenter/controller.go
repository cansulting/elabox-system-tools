// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// controls the initialization of services

package servicecenter

import (
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/web"
	"github.com/cansulting/elabox-system-tools/server"
)

/*
	serviceCenter.go
	file that manages the service connection
*/

////////////////////////GLOBAL DEFINITIONS///////////////

// initiate services
func Initialize(commandline bool) error {
	if commandline {
		return nil
	}
	global.Server = &server.Manager{}
	global.Server.OnSystemEvent = OnRecievedRequest
	global.Server.Setup()
	if err := global.Server.ListenAndServe(); err != nil {
		return err
	}
	webservice := &web.WebService{}
	if err := webservice.Start(); err != nil {
		return err
	}
	// start running all services
	return nil
}

/// this closes the server
func Close() {
	if global.Server != nil {
		global.Server.Stop()
	}
}
