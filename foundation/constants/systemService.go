package constants

const SYSTEM_SERVICE_ID = "ela.system"

// ------------------------ service center related Requests ---------------------
/*
app/service/action listener requests to listen to specific action.
They should subscribe to specific action via service center.
Check Service Center for references.
*/
const ACTION_SUBSCRIBE = "ela.system.action.SUBSCRIBE"

/*
app/service/action will broadcast an action to all listening app via service center.
Check Service Center for references.
*/
const ACTION_BROADCAST = "ela.system.action.BROADCAST"

/*
registers an application package so it will be available on the system.
Check Service Center for references.
*/
const SERVICE_REGISTER = "ela.system.action.REGISTER"

/*
unregisters a service. After the service is removed, it will not be executed upon
system initialization.
Check Service Center for references.
*/
const SERVICE_UNREGISTER = "ela.system.action.UNREGISTER"

/*
service state was changed. usually contains the state integer value.
Check Service Center for references.
*/
const SERVICE_CHANGE_STATE = "ela.system.action.SERVICE_STATE"

/*
App state changed. Check ApplicationState enum for possible values
*/
const APP_CHANGE_STATE = "ela.system.action.APP_STATE"
const SERVICE_PENDING_ACTIONS = "ela.system.service.PERNDING_ACTIONS"

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
