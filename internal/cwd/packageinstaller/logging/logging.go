package logging

import (
	"ela/foundation/path"
	"log"
	"os"
)

//var logString string = ""

func Initialize(filename string) error {
	return nil
	logPath := path.GetCacheDir() + "/" + filename
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
	/*
		for _, arg := range args {
			logString += string(arg)
		}*/
}

/*
func save(filename string) error {
	logPath := path.GetCacheDir() + "/logs/" + filename
	if err := os.MkdirAll(logPath, 0666); err != nil {
		return err
	}
	if err := os.WriteFile(logPath, []byte(logString), 0666); err != nil {
		return err
	}
	return nil
}
*/
