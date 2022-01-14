package main

import "github.com/cansulting/elabox-system-tools/foundation/app"

var controller *app.Controller

func main() {
	var err error
	controller, err = app.NewController(&MyActivity{}, &MyService{})
	if err != nil {
		panic(err)
	}

	app.RunApp(controller)
}
