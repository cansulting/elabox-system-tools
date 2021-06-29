package main

import (
	foundation "ela/foundation/app"
	"ela/foundation/app/data"
)

func main() {
	foundation.RunApp(&appmanager{}, data.AppData{Id: "ela.appmanager"})
}
