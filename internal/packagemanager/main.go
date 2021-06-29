package main

import (
	"ela/foundation/app"
	"ela/foundation/app/data"
)

func main() {
	app.RunApp(&mainApp{},
		data.AppData{Id: "ela.package-manager"})
}
