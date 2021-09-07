// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file contains definition for system action ids

package constants

//const ACTION_APP_AWAKE = "ela.action.APP_AWAKE"       // called when an app start running
//const ACTION_APP_SLEEP = "ela.action.APP_SLEEP"       // called when app stop running
const APP_CONFIG_NAME = "info.json"                           // default name for app config json
const ACTION_APP_LAUNCH = "ela.action.APP_LAUNCH"             // called when app will be launched
const ACTION_OPEN_CONTENT = "ela.action.CONTENT_OPEN"         // called when needs to open a content
const ACTION_GET_PENDING = "ela.action.GET_PENDING"           // called to retrieve pending actions for specific package
const ACTION_APP_INSTALL = "ela.action.INSTALL"               // called to launch app installer
const ACTION_APP_SYSTEM_INSTALL = "ela.action.SYSTEM_INSTALL" // called to initiate system installation
