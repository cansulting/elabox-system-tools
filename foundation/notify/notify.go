package notify

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

const PUSH_NOTIFICATION = constants.NOTIFICATION_CENTER_ID + ".action.PUSH_NOTIFICATION"

func System(val NotificationData) error {
	rpcInst, err := rpc.GetInstance()
	if err != nil {
		return err
	}
	notifData := data.NewAction(
		PUSH_NOTIFICATION,
		constants.NOTIFICATION_CENTER_ID,
		val,
	)
	if _, err := rpcInst.CallSystem(notifData); err != nil {
		return err
	}
	return nil
}
