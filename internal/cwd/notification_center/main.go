package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
)

func main() {

	con, err := app.NewController(nil, &MyService{})
	if err != nil {
		panic(err)
	}
	RPC = con.RPC
	AppController = con
	app.RunApp(con)
}
