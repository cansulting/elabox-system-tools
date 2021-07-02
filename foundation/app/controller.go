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
	event "ela/foundation/event"
	"ela/foundation/event/data"
	eventp "ela/foundation/event/protocol"
	"log"
	"time"
)

///////////////////////// FUNCTIONS ////////////////////////////////////

func RunApp(app *Controller) error {
	// start the app
	app.OnStart()
	log.Println("App is now running")

	for app.IsRunning() {
		time.Sleep(time.Second * 1)
	}

	log.Println("App exit")
	return app.OnEnd()
}

//////////////////////// CONTROLLER DEFINITION /////////////////////////////
// constructor for controller
func NewController(
	activity protocol.ActivityInterface,
	service protocol.ServiceInterface) (*Controller, error) {
	pk := appd.DefaultPackage()
	if err := pk.LoadFromSrc(constants.APP_CONFIG_NAME); err != nil {
		return nil, err
	}
	return &Controller{
		Service:  service,
		Activity: activity,
		Config:   pk,
	}, nil
}

type Controller struct {
	Service       protocol.ServiceInterface
	Activity      protocol.ActivityInterface
	SystemService *service.Connection // connection to system service
	Connector     eventp.ConnectorClient
	Config        *appd.PackageConfig
	running       bool
}

// true if this app is running
func (m *Controller) IsRunning() bool {
	return true
}

// callback when this app was started
func (m *Controller) OnStart() error {
	// step: init connector
	m.Connector = event.CreateClientConnector()
	err := m.Connector.Open()
	if err != nil {
		return err
	}
	// step: create service center connection
	m.SystemService, err = service.NewConnection(m.Connector, constants.SYSTEM_SERVICE_ID, onSystemServiceResponse)
	if err != nil {
		return err
	}
	// step: send running state
	_, err = m.SystemService.RequestFor(
		data.Action{
			Id:        constants.APP_CHANGE_STATE,
			PackageId: m.Config.PackageId,
			Data:      constants.APP_AWAKE})
	if err != nil {
		return err
	}
	// step: start service and activity
	if m.Service != nil {
		if err := m.Service.OnStart(); err != nil {
			return err
		}
	}
	if m.Activity != nil {
		if err := m.Activity.OnStart(); err != nil {
			return err
		}
	}
	return nil
}

// callback when this app ended
func (m *Controller) OnEnd() error {
	// step: send stop state for application
	_, err := m.SystemService.RequestFor(
		data.Action{
			Id:        constants.APP_CHANGE_STATE,
			PackageId: m.Config.PackageId,
			Data:      constants.APP_SLEEP})
	if err != nil {
		return err
	}
	if m.Activity != nil && m.Activity.IsRunning() {
		if err := m.Activity.OnEnd(); err != nil {
			return err
		}
	}
	if m.Service != nil && m.Service.IsRunning() {
		if err := m.Service.OnEnd(); err != nil {
			return err
		}
	}

	return nil
}

func onSystemServiceResponse(msg string, data interface{}) {

}
