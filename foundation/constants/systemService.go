package constants

const SYSTEM_SERVICE_ID = "ela.system"

// ------------------------ service center related Requests ---------------------
/*
app/service/action listener requests to listen to specific action.
They should subscribe to specific action via service center.
Check Service Center for references.
*/
const ACTION_SUBSCRIBE = "ela.system.SUBSCRIBE"

/*
app/service/action will broadcast an action to all listening app via service center.
Check Service Center for references.
*/
const ACTION_BROADCAST = "ela.system.BROADCAST"

/*
app/service/action will start an activity.
*/
const ACTION_START_ACTIVITY = "ela.system.START_ACTIVITY"

// system will be in update mode. all system will be terminated and only commandline will be available
const SYSTEM_UPDATE_MODE = "ela.system.UPDATE"

/*
servive is ready to serve
*/
const SERVICE_OPEN = "ela.system.service.OPEN"

/*
unregisters a service. After the service is removed, it will not be executed upon
system initialization.
Check Service Center for references.
*/
const SERVICE_CLOSE = "ela.system.service.CLOSE"

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
