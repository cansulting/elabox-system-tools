package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

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
	logger.InitFromFile("ela.testing", LOGFILE)
	logger.ClearLog()
	t.Log("Testing empty log...")
	src := LOGFILE
	reader, err := logger.NewReader(src)
	if err != nil {
		t.Error(err)
		return
	}
	reader.Load(0, 0, logger.LATEST_FIRST, nil)

	t.Log("Testing thousand of log...")

	logger.ConsoleOut = false
	for i := 0; i < 500; i++ {
		logger.GetInstance().Debug().Str("category", "testing").Int("index", i).Msg("This is testing number " + strconv.Itoa(i))
	}

	if err := getAllUnloadedLogs("debug", 500, reader); err != nil {
		t.Error(err)
		return
	}
	t.Log("Success!")
}

func TestLoggerQuery(t *testing.T) {
	t.Log("Testing Query Started...")
	logger.InitFromFile("ela.testing", LOGFILE)
	logger.ClearLog()

	// create mix log for debug, info and error
	var waitg sync.WaitGroup
	waitg.Add(3)
	errT := 100
	debugT := 150
	infoT := 200
	go func() {
		for i := 0; i < debugT; i++ {
			logger.GetInstance().Debug().Str("category", "testing").Int("index", i).Msg(strconv.Itoa(i) + " This is Debug")
		}
		waitg.Done()
	}()
	go func() {
		for i := 0; i < infoT; i++ {
			logger.GetInstance().Info().Str("category", "testing").Int("index", i).Msg(strconv.Itoa(i) + " This is Info")
		}
		waitg.Done()
	}()
	go func() {
		for i := 0; i < errT; i++ {
			logger.GetInstance().Error().Str("category", "testing").Int("index", i).Msg(strconv.Itoa(i) + " This is Error")
		}
		waitg.Done()
	}()
	waitg.Wait()

	reader, err := logger.NewReader(LOGFILE)
	if err != nil {
		t.Error(err)
		return
	}
	var errV int32 = 0
	var info int32 = 0
	var debug int32 = 0

	reader.Load(0, -1, logger.LATEST_FIRST, func(i int, l logger.Log) bool {
		switch l["level"] {
		case "error":
			atomic.AddInt32(&errV, 1)
		case "info":
			atomic.AddInt32(&info, 1)
		}
		return true
	})
	reader.Load(0, 0, logger.LATEST_FIRST, func(i int, l logger.Log) bool {
		switch l["level"] {
		case "debug":
			debug++
		}
		return true
	})
	if err := getAllUnloadedLogs("error", errT, reader); err != nil {
		t.Error(err)
		return
	}
	if err := getAllUnloadedLogs("info", infoT, reader); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Error =", errV, "Info =", info)
	if int(errV) != errT || int(info) != infoT {
		t.Error("Theres a missing log...")
		return
	}

	fmt.Println("Checking search query...")
	found := false
	var index2Search float64 = 0
	reader.Load(0, -1, logger.LATEST_FIRST, func(i int, l logger.Log) bool {
		if l["level"] == "error" && l["index"] == index2Search {
			found = true
			return false
		}
		return true
	})
	if !found {
		t.Error("Log with", index2Search, "not found")
		return
	} else {
		fmt.Println("Found search index", index2Search, "!")
	}

	fmt.Println("Check retrieve log with limit")
	retrieved, _ := reader.LoadSeq(-1, debugT, logger.LATEST_FIRST, func(l logger.Log) bool {
		return l["level"] == "debug"
	})
	if retrieved == debugT {
		fmt.Println("retrieve log with limit is working!")
	} else {
		t.Error("LoadLimit is not working")
	}

	fmt.Println("Check log with start from old")
	var lastRead float64 = 0
	retrieved, _ = reader.LoadSeq(-1, debugT, logger.OLD_FIRST, func(l logger.Log) bool {
		if l["level"] == "debug" {
			if l["index"] != lastRead {
				t.Error("Failed reading log starting from old")
				return false
			}
			lastRead++
		}
		return true
	})
	if retrieved == debugT {
		fmt.Println("retrieve log with limit is working!")
	} else {
		t.Error("Failed reading log starting from old")
	}
	t.Log("Success!")
}

func getAllUnloadedLogs(level string, totalLogs int, reader *logger.Reader) error {
	retrieved := make([]bool, totalLogs)
	founderror := ""
	reader.Load(0, -1, logger.LATEST_FIRST, func(i int, l logger.Log) bool {
		if l["level"] == level {
			index := int(l["index"].(float64))
			if !retrieved[index] {
				retrieved[index] = true
			} else {
				founderror = "Redundant value" + strconv.Itoa(index)
				return false
			}
		}
		return true
	})
	if founderror != "" {
		return errors.SystemNew(founderror, nil)
	}

	for i := 0; i < totalLogs; i++ {
		if !retrieved[i] {
			fmt.Println()
			return errors.SystemNew("Unretrieved index"+strconv.Itoa(i), nil)
		}
	}

	fmt.Println("All", level, "was retrieved")
	return nil
}
