//go:build IDE && !RELEASE && !STAGING
// +build IDE,!RELEASE,!STAGING

package global

const ENV = "debug"
const SYSVER_HOST = "https://storage.googleapis.com/elabox-debug/packages"