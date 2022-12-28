package main

import "github.com/cansulting/elabox-system-tools/foundation/event/data"

func BroadcastNotification(notification NotifData) error {
	_, err := RPC.CallBroadcast(data.NewAction(BROADCAST_NOTIFICATION, "", notification))
	return err
}
