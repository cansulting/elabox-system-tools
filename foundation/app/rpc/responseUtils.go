// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// this file provides funnction for generating response data for RPC
// response data is json string which will be encoded to Base64 before transmitting

package rpc

import (
	"encoding/base64"
	"strconv"
	"strings"
)

const SUCCESS_CODE = 200
const SYSTEMERR_CODE = 400 // theres something wrong with the system
const INVALID_CODE = 401
const NOT_IMPLEMENTED = 300 // code was not implemented

// return json string for response
func CreateResponse(code int16, msg string) string {
	return CreateResponseQ(code, msg, true)
}

func CreateResponseQ(code int16, msg string, addQoute bool) string {
	if addQoute {
		if msg != "" {
			msg = strings.Replace(msg, "\"", "\\\"", -1)
		}
		msg = "\"" + msg + "\""
	}
	return base64.StdEncoding.EncodeToString([]byte("{\"code\":" + strconv.Itoa(int(code)) + ", \"message\": " + msg + "}"))
}

func CreateJsonResponse(code int16, msg string) string {
	return CreateResponseQ(code, msg, false)
}

// returns success json response
func CreateSuccessResponse(msg string) string {
	return CreateResponseQ(SUCCESS_CODE, msg, true)
}

func DecodeResponse(base64str string) (string, error) {
	base64str = strings.Replace(base64str, "\"", "", -1)
	res, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
