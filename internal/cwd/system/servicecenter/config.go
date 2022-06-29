package servicecenter

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/internal/cwd/env"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/appman"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
)

func configureSystem() string {
	if !isConfig() {
		if err := env.SetEnv(global.CONFIG_ENV, "1"); err != nil {
			global.Logger.Error().Err(err).Msg("failed to mark as config")
		}
		appman.InitializeAllPackages()
	}
	return rpc.CreateSuccessResponse("success")
}

func isConfig() bool {
	return env.GetEnv(global.CONFIG_ENV) == "1"
}
