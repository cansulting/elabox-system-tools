package constants

const SERVICE_CENTER_ID = "ela.service-center"

// ------------------------ service center related Requests ---------------------
/*
app/service/action listener requests to listen to specific action.
They should subscribe to specific action via service center.
Check Service Center for references.
*/
const ACTION_SUBSCRIBE = "ela.service-center.action.SUBSCRIBE"

/*
app/service/action will broadcast an action to all listening app via service center.
Check Service Center for references.
*/
const ACTION_BROADCAST = "ela.service-center.action.BROADCAST"

/*
registers a service. Service needs to be registered before become fully usable.
It will executed upon system initialization. Check Service Center for references.
*/
const SERVICE_REGISTER = "ela.service-center.action.REGISTER"

/*
unregisters a service. After the service is removed, it will not be executed upon
system initialization.
Check Service Center for references.
*/
const SERVICE_UNREGISTER = "ela.service-center.action.UNREGISTER"

/*
service state was changed. usually contains the state integer value.
Check Service Center for references.
*/
const SERVICE_CHANGE_STATE = "ela.service-center.action.CHANGE_STATE"

type ServiceState int

const (
	SERVICE_AWAKE = 1 // service is currently running
	SERVICE_SLEEP = 0 // service is not yet executed
)
