// handles information for app state

package data

import "github.com/cansulting/elabox-system-tools/foundation/constants"

type AppState struct {
	State constants.AppRunningState `json:"state"`
	Data  interface{}               `json:"data"`
}
