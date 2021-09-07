// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package constants

const SYSTEM_SERVICE_ID = "ela.system"

// ------------------------ service center related Requests ---------------------
// app/service/action listener requests to listen to specific action.
// They should subscribe to specific action via service center.
// Check Service Center for references.
const ACTION_SUBSCRIBE = "ela.system.SUBSCRIBE"
const ACTION_SYSTEM_STATUS = "ela.system.STATUS"

/*
app/service/action will broadcast an action to all listening app via service center.
Check Service Center for references.
*/
const SYSTEM_BROADCAST = "ela.system.BROADCAST"

/*
app/service/action will start an activity.
*/
const ACTION_START_ACTIVITY = "ela.system.START_ACTIVITY"

// system will be in update mode. all system will be terminated and only commandline will be available
const SYSTEM_UPDATE_MODE = "ela.system.UPDATE"

// called when an activity returns a result
const SYSTEM_ACTIVITY_RESULT = "ela.system.ACTIVITY_RESULT"

// system will terminate
const SYSTEM_TERMINATE = "ela.system.TERMINATE"
const SYSTEM_TERMINATE_NOW = "ela.system.TERMINATE_NOW"

/*
service state was changed. usually contains the state integer value.
Check Service Center for references.
*/
const SERVICE_CHANGE_STATE = "ela.system.SERVICE_STATE"

/*
App state changed. Check ApplicationState enum for possible values
*/
const APP_CHANGE_STATE = "ela.system.action.APP_STATE"
const SERVICE_PENDING_ACTIONS = "ela.system.PENDING_ACTIONS"

// sends terminate action to app
const APP_TERMINATE = "ela.system.APP_TERMINATE"

type ServiceState int

const (
	SERVICE_AWAKE = 1 // service is currently running
	SERVICE_SLEEP = 0 // service is not yet executed
)

type ApplicationState int

const (
	APP_AWAKE = 1 // app is running
	APP_SLEEP     // app was stopped
)
