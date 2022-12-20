package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
)

const PACKAGE_ID = "ela.notification"
const NOTIF_QUEUE_LIMIT = 100 // maximum content for notification queue

// ACTIONS
const AC_PUSH_NOTIF = PACKAGE_ID + ".action.PUSH_NOTIFICATION"
const AC_RETRIEVE_NOTIF = PACKAGE_ID + ".action.RETRIEVE_NOTIF"

// BROADCASTS
const BROADCAST_NOTIFICATION = PACKAGE_ID + ".broadcast.PUSH_NOTIFICATION"

var AppController *app.Controller
var RPC *rpc.RPCHandler
