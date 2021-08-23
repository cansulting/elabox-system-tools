// +build !RELEASE,!STAGING

package config

func GetBuildMode() Mode {
	return DEBUG
}
