// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package constants

import "github.com/cansulting/elabox-system-tools/foundation/path"

// CONNECTION CONFIG
const PORT = 80
const TIMEOUT = 20        // transmit timeout
const RECONNECT = true    // true if enable reconnection when reconnecting to socket io
const PACKAGE_EXT = "box" // extension of packaged file
const CUSTOMINSLLER = "custominstaller"

const LOG_FILE = path.PATH_LOG + "/elabox.log" // the system log file path

// package ids
const SYSTEM_SERVICE_ID = "ela.system"
const NOTIFICATION_CENTER_ID = "ela.notification"
const ACCOUNT_SYS_ID = "ela.account"
