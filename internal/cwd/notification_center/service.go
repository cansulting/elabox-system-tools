package main

import (
	"fmt"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
)

type MyService struct {
}

func (instance *MyService) OnStart() error {
	// register service rpc
	AppController.RPC.OnRecieved(AC_PUSH_NOTIF, instance.onPushNotification)
	AppController.RPC.OnRecieved(AC_RETRIEVE_NOTIF, instance.onRetrieveNotification)
	//go RetrieveAllApps(false)
	return nil
}

func (instance *MyService) onPushNotification(client protocol.ClientInterface, action data.Action) string {
	data := NotifData{}
	if err := action.DataToObj(&data); err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	// no icon provided? add package icon
	if data.Icon == "" {
		data.Icon = "https://365webresources.com/wp-content/uploads/2013/11/iOS-7-Icon-Grid.png"
	}
	data.ReceivedAt = time.Now().UTC()
	data.Id = fmt.Sprintf("%d", data.ReceivedAt.Unix())
	if err := AddNotif(data); err != nil {
		return rpc.CreateResponse(rpc.INVALID_PARAMETER_PROVIDED, err.Error())
	}
	// Broadcast notification
	if err := BroadcastNotification(data); err != nil {
		logger.GetInstance().Error().Err(err).Msg("failed to broacast push notification event")
	}
	return rpc.CreateSuccessResponse("success")
}

func (instance *MyService) onRetrieveNotification(client protocol.ClientInterface, action data.Action) string {
	params, err := action.DataToMap()
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	length := 10
	page := 1
	if params["page"] != nil {
		page = int(params["page"].(float64))
	}
	if params["length"] != nil {
		length = int(params["length"].(float64))
	}
	notifications, err := RetrieveNotif(uint(page), uint(length))
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	return rpc.CreateJsonResponse(rpc.SUCCESS_CODE, notifications)
}

func (instance *MyService) IsRunning() bool {
	return true
}

func (instance *MyService) OnEnd() error {
	return nil
}
