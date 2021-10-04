package main

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Log("Check the log file @" + constants.LOG_FILE)
	// we need to initialize logger by passing the current app package
	logger.Init("ela.testing")

	// for debug level. Debug is for additional information specifically for debugging.
	logger.GetInstance().Debug().Msg("Hello")

	// you can also add category of log for any level(debug, error, info, warn, fatal, panic)
	// category can be use to categorize the logging. which can be filter later
	// example below is for networking category
	logger.GetInstance().Debug().Str("category", "networking").Msg("This is a sample debug with networking category")

	// for error level. error log includes error data for thrown errors
	var sampleError error = errors.SystemNew("Sample error", nil)
	logger.GetInstance().Error().Err(sampleError).Msg("This is an error")
	// you can add error with caller. caller includes caller function and line number
	logger.GetInstance().Error().Err(sampleError).Caller().Msg("This is an error with sample caller")
	// you can add stackstrace. this log the call stack
	logger.GetInstance().Error().Err(sampleError).Stack().Msg("This is an error with sample stack")

	// for warning. for logs that warns
	logger.GetInstance().Warn().Str("category", "system").Msg("This is a sample warning with category")
	
	// for fatal and panic related
	//logger.GetInstance().Fatal().Stack().Msg("This is fatal log")
	//logger.GetInstance().Panic().Stack().Msg("This is panic log")

	t.Log("Sucess")
}
