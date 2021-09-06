package config

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/foundation/perm"
	"encoding/json"
	"log"
	"os"
)

/*
	Structure for handling system environment variables.
	This structure includes serialization and deserialization.
	TODO: optimize serialization. Currently this serializes everytime theres a changes with environment variable
*/
var singleton *Env

const FILENAME = "env.json"

type Env struct {
	Vars map[string]string `json:"vars"`
}

func getSourceLoc() string {
	return path.GetSystemAppDirData(constants.SYSTEM_SERVICE_ID) + "/" + FILENAME
}

func initialize() error {
	if singleton == nil {
		singleton = &Env{}
		envsrc := getSourceLoc()
		if _, err := os.Stat(envsrc); err == nil {
			bytes, err := os.ReadFile(envsrc)
			if err != nil {
				return errors.SystemNew("Failed loading environment file @ "+envsrc+" ", err)
			}
			if len(bytes) > 0 {
				if err := json.Unmarshal(bytes, &singleton); err != nil {
					return errors.SystemNew("Failed loading environment file @ "+envsrc+" ", err)
				}
			}
		}
		if singleton.Vars == nil {
			singleton.Vars = map[string]string{}
		}
	}
	return nil
}

func SetEnv(key string, value string) error {
	if err := initialize(); err != nil {
		return err
	}
	singleton.Vars[key] = value
	os.Setenv(key, value)
	return SaveEnv()
}

func GetEnv(key string) string {
	if err := initialize(); err != nil {
		log.Panic(err)
		return ""
	}
	if key == "" {
		return ""
	}
	return singleton.Vars[key]
}

func SaveEnv() error {
	marshaled, err := json.Marshal(singleton)
	if err != nil {
		return errors.SystemNew("Failed saving environment config file", err)
	}
	if err := os.WriteFile(getSourceLoc(), marshaled, perm.PUBLIC_VIEW); err != nil {
		return errors.SystemNew("Failed saving environment config file", err)
	}
	return nil
}
