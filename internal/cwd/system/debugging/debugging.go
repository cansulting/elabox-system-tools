// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file provides extension functions for debugging

package debugging

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/appman"
)

// use to debug a package. if package is already running then stop it, hence create a debug app
func DebugApp(pkid string, client protocol.ClientInterface) *appman.AppConnect {
	app := appman.LookupAppConnect(pkid)
	if app != nil && app.Client != nil && app.Client.IsAlive() {
		return app
	}
	// step: stop if app is already running
	appman.RemoveAppConnect(pkid, true)
	return appman.AddAppConnect(&data.PackageConfig{PackageId: pkid}, client)
}
