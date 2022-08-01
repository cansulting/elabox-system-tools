package appman

import (
	"errors"

	"github.com/cansulting/elabox-system-tools/foundation/system"
)

func GetDeviceSerial() (string, error) {
	deviceSerial := system.GetDeviceInfo().Serial
	if len(deviceSerial) == 0 {
		return "", errors.New("Unable to get device serial")
	}
	return deviceSerial, nil
}
