// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package servicecenter

import "strconv"

const SUCCESS_CODE = 200
const SYSTEMERR_CODE = 400 // theres something wrong with the system
const INVALID_CODE = 401

// return json string for response
func CreateResponse(code int16, msg string) string {
	return `{"code":` + strconv.Itoa(int(code)) + `, "message": "` + msg + `"}`
}

// returns success json response
func CreateSuccessResponse(msg string) string {
	return CreateResponse(SUCCESS_CODE, msg)
}
