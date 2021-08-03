package main

import (
	"ela/foundation/app"
	"log"
)

var controller *app.Controller

func main() {
	_controller, err := app.NewController(nil, &MyService{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	controller = _controller
	if err != nil {
		log.Fatal(err.Error())
	}
	app.RunApp(controller)
	/*
		s := MyService{}
		if err := s.OnStart(); err != nil {
			log.Println(err)
			return
		}
		for s.IsRunning() {
			go time.Sleep(time.Second * 2)
		}*/
}
