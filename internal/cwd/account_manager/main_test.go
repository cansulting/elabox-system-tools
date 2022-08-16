package main

import (
	"testing"

	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/factory"
)

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

func Test_AccountCreate(t *testing.T) {
	acc := "test"
	if exist, _ := factory.IsAccountExist(acc); exist {
		if err := factory.DeleteAccount(acc); err != nil {
			t.Error("failed deleting account", acc, err)
			return
		}
	}
	acc_data := data.Account{
		Username: acc,
	}
	if err := factory.CreateAccount(acc, acc_data); err != nil {
		t.Error("failed creating account", err)
	}
}
