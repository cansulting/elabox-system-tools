package main

import "testing"

const SAMPLE_DID = "HELLO_DID_USER"
const SAMPLE_USERNAME = "elabox"
const SAMPLE_PASSWORD = "elabox"

// use to test device did
func Test_SetDeviceDid(t *testing.T) {
	presentation := make(map[string]interface{})
	presentation["holder"] = SAMPLE_DID
	if err := SetDeviceDid(presentation); err != nil {
		t.Error(err)
	}
}

func Test_AuthenticationDid(t *testing.T) {
	if !AuthenticateDid(SAMPLE_DID) {
		t.Error("failed to authenticate did")
	}
}

func Test_AuthenticateSystemAccount(t *testing.T) {
	err, success := AuthenticateSystemAccount(SAMPLE_USERNAME, SAMPLE_PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Authenticate", success)
}
func Test_GetDeviceSerial(t *testing.T) {
	deviceSerial := GetDeviceSerial()
	if len(deviceSerial) == 0 {
		t.Error("Error in getting device serial")
		return
	}
	t.Log("Got Device Serial", deviceSerial)
}
