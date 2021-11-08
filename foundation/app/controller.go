package app

// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// controller.go
// Controller class handles the application lifecycle.
// To initialize call NewController, for debugging use NewControllerWithDebug
// please see the documentation for more info.
import (
	"strconv"
	"time"

	appd "github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/app/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	event "github.com/cansulting/elabox-system-tools/foundation/event"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	protocolE "github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/system"
)

///////////////////////// FUNCTIONS ////////////////////////////////////

func RunApp(app *Controller) error {
	// start the app
	if err := app.onStart(); err != nil {
		return err
	}
	logger.GetInstance().Info().Str("category", "appcontroller").Msg(app.Config.PackageId + "is now running")

	for app.IsRunning() {
		time.Sleep(time.Second * 1)
	}

	defer logger.GetInstance().Info().Str("category", "appcontroller").Msg("App exit " + app.Config.PackageId)
	return app.onEnd()
}

//////////////////////// CONTROLLER DEFINITION /////////////////////////////
// constructor for controller
// @activity the activity function for this app
// @service the service function for this app
// Please check the system app manager debugapp()
func NewController(
	activity protocol.ActivityInterface,
	service protocol.ServiceInterface) (*Controller, error) {

	config := appd.DefaultPackage()
	if err := config.LoadFromSrc(constants.APP_CONFIG_NAME); err != nil {
		if logger.GetInstance() == nil {
			logger.Init(config.PackageId)
		}
		logger.GetInstance().Panic().Err(err).Msg("Unable to find package info.json")
	}
	if logger.GetInstance() == nil {
		logger.Init(config.PackageId)
	}
	return &Controller{
		Debugging:  system.IDE,
		AppService: service,
		Activity:   activity,
		Config:     config,
	}, nil
}

type Controller struct {
	AppService protocol.ServiceInterface // current service for this app
	Activity   protocol.ActivityInterface
	RPC        rpc.RPCInterface //
	Config     *appd.PackageConfig
	forceEnd   bool
	Debugging  bool // true if the app currently debugging
}

// true if this app is running
func (m *Controller) IsRunning() bool {
	if m.forceEnd {
		return false
	}
	if m.Activity != nil && m.Activity.IsRunning() {
		return true
	}
	if m.AppService != nil && m.AppService.IsRunning() {
		return true
	}
	return false
}

// callback when this app was started
func (m *Controller) onStart() error {
	logger.GetInstance().
		Info().
		Str("category", "appcontroller").
		Msg("Starting App " + m.Config.PackageId + ". Ide = " + strconv.FormatBool(system.IDE))
	// step: init connector
	connector := event.CreateClientConnector()
	err := connector.Open(-1)
	if err != nil {
		return errors.SystemNew("Controller: Failed to start. Couldnt create client connector.", err)
	}
	// step: create RPC
	if m.RPC == nil {
		m.RPC = rpc.NewRPCHandler(connector)
		m.RPC.OnRecieved(constants.APP_TERMINATE, m.onTerminate)
	}
	// step: send running state
	awake := constants.APP_AWAKE
	if m.Debugging {
		awake = constants.APP_AWAKE_DEBUG
	}
	res, err := m.RPC.CallSystem(
		data.NewAction(constants.APP_CHANGE_STATE, m.Config.PackageId, awake))
	if err != nil {
		return err
	}
	println(m.Config.PackageId + "onstart pendingActions =" + res.ToString())
	pendingActions := res.ToActionGroup()
	// step: initialize service
	if m.AppService != nil {
		logger.GetInstance().Debug().Str("category", "appcontroller").Msg("Service start")
		if err := m.AppService.OnStart(); err != nil {
			return errors.SystemNew("app.Controller couldnt start app service", err)
		}
	}
	// step: initialize activity
	if m.Activity != nil {
		logger.GetInstance().Debug().Str("category", "appcontroller").Msg("Activity start")
		if err := m.Activity.OnStart(pendingActions.Activity); err != nil {
			return errors.SystemNew("app.Controller couldnt start app activity", err)
		}
	}
	return nil
}

// callback when this app ended
func (m *Controller) onEnd() error {
	//log.Println("Controller: OnEnd")
	if m.forceEnd {
		// step: send stop state for application
		_, err := m.RPC.CallSystem(
			data.NewAction(
				constants.APP_CHANGE_STATE,
				m.Config.PackageId,
				constants.APP_SLEEP))
		if err != nil {
			logger.GetInstance().Error().Err(err).Caller().Str("category", "appcontroller").Msg("Controller.onEnd Change state failed.")
		}
	}
	if m.Activity != nil && m.Activity.IsRunning() {
		if err := m.Activity.OnEnd(); err != nil {
			logger.GetInstance().Error().Err(err).Caller().Str("category", "appcontroller").Msg("Activity stop failed")
		}
	}
	if m.AppService != nil && m.AppService.IsRunning() {
		if err := m.AppService.OnEnd(); err != nil {
			logger.GetInstance().Error().Err(err).Caller().Str("category", "appcontroller").Msg("AppService stop failed")
		}
	}
	return nil
}

// this will end the app
func (c *Controller) End() {
	logger.GetInstance().Debug().Str("category", "appcontroller").Msg(c.Config.PackageId + "is now ending")
	c.forceEnd = true
}

// use to start an  from other applications
func (m *Controller) StartActivity(action data.Action) error {
	logger.GetInstance().Debug().Str("category", "appcontroller").Msg("Trying to start activity with action" + action.Id)
	res, err := m.RPC.CallSystem(data.NewAction(constants.ACTION_START_ACTIVITY, "", action))
	if err != nil {
		return err
	}
	logger.GetInstance().Debug().Str("category", "appcontroller").Msg("Start activity with response " + res.ToString())
	return nil
}

// use to return result to the caller of this app
func (c *Controller) SetActivityResult(val interface{}) {
	res, err := c.RPC.CallSystem(data.NewAction(constants.SYSTEM_ACTIVITY_RESULT, c.Config.PackageId, val))
	if err != nil {
		logger.GetInstance().Error().Str("category", "appcontroller").Err(err).Caller().Msg("Activity result response failure")
		return
	}
	if res != nil {
		logger.GetInstance().Debug().Str("category", "appcontroller").Msg("Activity result response success " + res.ToString())
	}
}

// callback from system. this app will be terminated
func (c *Controller) onTerminate(client protocolE.ClientInterface, data data.Action) string {
	c.End()
	return ""
}
