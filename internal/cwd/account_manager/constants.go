package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

// actions
const AC_AUTH_DID = "account.actions.AUTH_DID"                   // use to authenticaticate specific did
const AC_SETUP_CHECK = "account.actions.DID_SETUP_CHECK"         // use to check if theres an existing did setup
const AC_SETUP_DID = "account.actions.DID_SETUP"                 // use to setup did
const AC_SETUP_ACCOUNT = "account.actions.ACC_SETUP"             // use to setup user account
const AC_VALIDATE_TOKEN = PACKAGE_ID + ".actions.VALIDATE_TOKEN" // use to validate token
const AC_AUTH_SYSTEM = PACKAGE_ID + ".actions.AUTH_SYSTEM"       // authenticate via system
const AC_PASS_CHANGE = PACKAGE_ID + ".actions.PASS_CHANGE"       // use to update password

const PACKAGE_ID = "ela.account"
const KEYSTORE_PATH = "/home/elabox/documents/ela.mainchain/keystore.dat"
const DEFAULT_USERNAME = "elabox"

var ACCOUNT_LOC = path.GetSystemAppDirData(PACKAGE_ID) + "/user"

var Controller *app.Controller

// error codes
const CREDENTIALS_INVALID = 800
