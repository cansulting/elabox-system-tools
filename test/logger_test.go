package main

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
)

const LOGFILE = "log.txt"

func TestLoggerWrite(t *testing.T) {
	// we need to initialize logger by passing the current app package
	//logger.Init("ela.testing")
	// or you can use use pre defined log file
	logger.InitFromFile("ela.testing", "log.txt")

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

func TestLoggerRead(t *testing.T) {
	t.Log("Testing empty log...")
	os.Remove(LOGFILE)
	src := LOGFILE
	reader, err := logger.NewReader(src)
	if err != nil {
		t.Error(err)
		return
	}
	reader.Load()

	t.Log("Testing thousand of log...")
	logger.InitFromFile("ela.testing", LOGFILE)
	for i := 0; i < 1000; i++ {
		logger.GetInstance().Debug().Str("category", "testing").Msg("This is testing number " + strconv.Itoa(i))
	}

	c := make(chan int)
	go func(c chan int) {
		time.Sleep(time.Second * 7)
		reader.Load()
		t.Log("Success!")
		c <- 1
	}(c)

	<-c
}
