package env

import (
	"encoding/json"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
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

func initEnv() error {
	if singleton == nil {
		singleton = &Env{}
		envsrc := getSourceLoc()
		if _, err := os.Stat(envsrc); err == nil {
			// load from file
			bytes, err := os.ReadFile(envsrc)
			if err != nil {
				return errors.SystemNew("Failed loading environment file @ "+envsrc+" ", err)
			}
			if len(bytes) > 0 {
				if err := json.Unmarshal(bytes, &singleton); err != nil {
					return errors.SystemNew("Failed loading environment file @ "+envsrc+" ", err)
				}
				// env loaded
				for key, val := range singleton.Vars {
					os.Setenv(key, val)
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
	if err := initEnv(); err != nil {
		return err
	}
	if singleton.Vars[key] != value {
		singleton.Vars[key] = value
		os.Setenv(key, value)
		return saveEnv()
	}
	return nil
}

func GetEnv(key string) string {
	if err := initEnv(); err != nil {
		return ""
	}
	if key == "" {
		return ""
	}
	return singleton.Vars[key]
}

func saveEnv() error {
	marshaled, err := json.Marshal(singleton)
	if err != nil {
		return errors.SystemNew("Failed saving environment config file", err)
	}
	if err := os.WriteFile(getSourceLoc(), marshaled, perm.PUBLIC_VIEW); err != nil {
		return errors.SystemNew("Failed saving environment config file", err)
	}
	return nil
}
