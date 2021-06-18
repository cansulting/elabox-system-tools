package main

import (
	base "ela.services/Base"
)

func main() {
	base.RunApp(&ApplicationManager{}, base.AppData{Id: "ela.ApplicationManager"})
}
