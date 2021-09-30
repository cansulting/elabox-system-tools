// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package logger

import (
	"ela/foundation/constants"
	"ela/foundation/perm"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var instanceLogger *zerolog.Logger
var ConsoleOut = true // true if write log on console not only in file

// this creates a new log if not yet created
func Init(packageId string) *zerolog.Logger {
	//if instanceLogger == nil {
	// init logfile
	logfile, err := os.OpenFile(constants.LOG_FILE, os.O_CREATE|os.O_RDWR|os.O_APPEND, perm.PUBLIC_WRITE)
	if err != nil {
		fmt.Println("Error opening logfile "+constants.LOG_FILE, err)
		return nil
	}
	fmt.Println("Log file opened @", constants.LOG_FILE)
	var writer io.Writer = logfile
	if ConsoleOut {
		writer = zerolog.MultiLevelWriter(logfile, os.Stdout)
	}
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(writer).With().Timestamp().Str("package", packageId).Logger()
	instanceLogger = &logger
	return instanceLogger
}

// get the current instance of logger
func GetInstance() *zerolog.Logger {
	return instanceLogger
}

// use to set hook
func SetHook(h zerolog.Hook) {
	logger := instanceLogger.Hook(h)
	instanceLogger = &logger
}
