package main

import "github.com/cansulting/elabox-system-tools/foundation/system"

func GetDeviceSerial() string {
	deviceSerial := system.GetDeviceInfo().Serial
	return deviceSerial
}
