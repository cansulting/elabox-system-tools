package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
)

func main() {

	con, err := app.NewController(nil, &MyService{})
	if err != nil {
		panic(err)
	}
	global.RPC = con.RPC
	global.AppController = con
	app.RunApp(con)
}
