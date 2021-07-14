package app

/*
	controller.go
	Provides basic implementation application.
	Application can contains service, activity and broadcast listener
*/
import (
	appd "ela/foundation/app/data"
	"ela/foundation/app/protocol"
	"ela/foundation/app/service"
	"ela/foundation/constants"
	"ela/foundation/errors"
	event "ela/foundation/event"
	"ela/foundation/event/data"
	eventp "ela/foundation/event/protocol"
	"log"
	"time"
)

///////////////////////// FUNCTIONS ////////////////////////////////////

func RunApp(app *Controller) error {
	// start the app
	if err := app.OnStart(); err != nil {
		return err
	}
	log.Println("App is now running")

	for app.IsRunning() {
		time.Sleep(time.Second * 1)
	}

	defer log.Println("App exit")
	return app.OnEnd()
}

//////////////////////// CONTROLLER DEFINITION /////////////////////////////
// constructor for controller
func NewController(
	activity protocol.ActivityInterface,
	service protocol.ServiceInterface) (*Controller, error) {
	config := appd.DefaultPackage()
	if err := config.LoadFromSrc(constants.APP_CONFIG_NAME); err != nil {
		return nil, err
	}
	return &Controller{
		AppService: service,
		Activity:   activity,
		Config:     config,
	}, nil
}

type Controller struct {
	AppService    protocol.ServiceInterface // current service for this app
	Activity      protocol.ActivityInterface
	SystemService service.IConnection // connection to system service. system service handles system main operation like app registration
	Connector     eventp.ConnectorClient
	Config        *appd.PackageConfig
	forceEnd      bool
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
func (m *Controller) OnStart() error {
	log.Println("app.Controller: Starting App", m.Config.PackageId)
	// step: init connector
	m.Connector = event.CreateClientConnector()
	err := m.Connector.Open(-1)
	if err != nil {
		return errors.SystemNew("Controller: Failed to start. Couldnt create client connector.", err)
	}
	// step: create service center connection
	if m.SystemService == nil {
		m.SystemService, err = service.NewConnection(
			m.Connector,
			constants.SYSTEM_SERVICE_ID,
			m.onSystemServiceResponse,
		)
		if err != nil {
			return errors.SystemNew(
				"Controller: Failed to start. Couldnt connect to service "+constants.SYSTEM_SERVICE_ID, err)
		}
	}
	// step: send running state
	res, err := m.SystemService.RequestFor(
		data.NewAction(
			constants.APP_CHANGE_STATE,
			m.Config.PackageId,
			constants.APP_AWAKE))
	if err != nil {
		return err
	}
	log.Println("controller.OnStart() pendingActions =", res)
	pendingActions := res.ToActionGroup()
	// step: start service and activity
	if m.AppService != nil {
		log.Println("app.Controller: OnStart", "Service")
		if err := m.AppService.OnStart(); err != nil {
			return errors.SystemNew("app.Controller couldnt start app service", err)
		}
	}
	if m.Activity != nil {
		log.Println("app.Controller: OnStart", "Activity")
		if err := m.Activity.OnStart(pendingActions.Activity); err != nil {
			return errors.SystemNew("app.Controller couldnt start app activity", err)
		}
	}
	return nil
}

// callback when this app ended
func (m *Controller) OnEnd() error {
	log.Println("Controller: OnEnd")
	if m.forceEnd {
		// step: send stop state for application
		_, err := m.SystemService.RequestFor(
			data.NewAction(
				constants.APP_CHANGE_STATE,
				m.Config.PackageId,
				constants.APP_SLEEP))
		if err != nil {
			return err
		}
	}
	if m.Activity != nil && m.Activity.IsRunning() {
		if err := m.Activity.OnEnd(); err != nil {
			return err
		}
	}
	if m.AppService != nil && m.AppService.IsRunning() {
		if err := m.AppService.OnEnd(); err != nil {
			return err
		}
	}

	return nil
}

// this will end the app
func (c *Controller) End() {
	c.forceEnd = true
}

// use to start an  from other applications
func (m *Controller) StartActivity(action data.Action) error {
	log.Println("Controller:StartActivity", action.Id)
	res, err := m.SystemService.RequestFor(data.NewAction(constants.ACTION_START_ACTIVITY, "", action))
	if err != nil {
		return err
	}
	log.Println("Controller:StartActivity response", res.ToString())
	return nil
}

func (c *Controller) onSystemServiceResponse(msg string, data interface{}) {
	switch msg {
	case constants.APP_TERMINATE:
		c.End()
	}
}
