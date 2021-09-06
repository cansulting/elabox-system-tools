package config

import (
	"ela/foundation/errors"
)

const ELABUILD = "ELABUILD"					// environment variable for build mode

// intialize system configuration
func Init() error {
	if err := SetEnv(ELABUILD, string(GetBuildMode())); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	return nil
}
