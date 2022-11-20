//go:build IDE && !RELEASE && !STAGING
// +build IDE,!RELEASE,!STAGING

package global

const ENV = "debug"
const RETRIEVE_LISTING_DELAY = 5 // delay in retrieving store listing in sec
const STOREHUB_SERVER = "http://localhost:4005"
