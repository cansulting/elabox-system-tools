//go:build RELEASE
// +build RELEASE

package global

const ENV="release"
const RETRIEVE_LISTING_DELAY = 60 * 60 // delay in retrieving store listing in sec
const STOREHUB_SERVER = "http://localhost:4005"
