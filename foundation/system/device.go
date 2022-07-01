/// Copyright 2021 The Elabox Authors
/// This file is part of the elabox-rewards library.
///
/// NOTICE:  All information contained herein is, and remains the property of Elabox. The intellectual and technical concepts contained
/// herein are proprietary to Elabox and may be covered by U.S. and Foreign Patents, patents in process, and are protected by trade secret or copyright law.
/// Dissemination of this information or reproduction of this material is strictly forbidden unless prior written permission is obtained
/// from Elabox.  Access to the source code contained herein is hereby forbidden to anyone except current Elabox employees, managers or contractors who have executed
/// Confidentiality and Non-disclosure agreements explicitly covering such access.
///
/// The copyright notice above does not evidence any actual or intended publication or disclosure  of  this source code, which includes
/// information that is confidential and/or proprietary, and is a trade secret, of  Elabox.   ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC  PERFORMANCE,
/// OR PUBLIC DISPLAY OF OR THROUGH USE  OF THIS  SOURCE CODE  WITHOUT  THE EXPRESS WRITTEN CONSENT OF ELABOX IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE
/// LAWS AND INTERNATIONAL TREATIES.  THE RECEIPT OR POSSESSION OF  THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY OR IMPLY ANY RIGHTS
/// TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING THAT IT  MAY DESCRIBE, IN WHOLE OR IN PAR

package system

import (
	"os"
	"strings"
)

type DeviceInfo struct {
	Model    string
	Serial   string
	Hardware string
}

// use to retrieve current device information
func GetDeviceInfo() DeviceInfo {
	var device DeviceInfo
	fname := "/proc/cpuinfo"
	bytes, _ := os.ReadFile(fname)
	splits := strings.Split(string(bytes), "\n")
	for _, split := range splits {
		keypair := strings.Split(split, ":")
		if len(keypair) > 0 {
			key := strings.Trim(keypair[0], "\t")
			switch key {
			case "Serial":
				device.Serial = strings.Trim(keypair[1], " ")
			case "Hardware":
				device.Hardware = strings.Trim(keypair[1], " ")
			case "Model":
				device.Model = strings.Trim(keypair[1], " ")
			}
		}
	}
	return device
}
