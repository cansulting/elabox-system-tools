package base

const ACTIONCENTER_ID = "ela.service-center"
const ACTION_SUBSCRIBE = "ela.service-center.action.SUBSCRIBE"
const ACTION_BROADCAST = "ela.service-center.action.BROADCAST"
const ACTION_REGISTER = "ela.service-center.action.REGISTER"
const ACTION_UNREGISTER = "ela.service-center.action.UNREGISTER"
const ACTION_CHANGE_STATE = "ela.service-center.action.CHANGE_STATE"

type ServiceState int

const (
	ServiceAwake    = iota // service means currently running
	ServiceSleeping        // service is not executed
	ServiceStopped         // service was stopped
)
