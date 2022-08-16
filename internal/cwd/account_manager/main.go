package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/data"
)

func main() {
	con, err := app.NewController(nil, &MyService{})
	if err != nil {
		panic(err)
	}
	data.Controller = con
	app.RunApp(con)
}
