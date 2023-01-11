package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

// actions
const AC_AUTH_DID = "account.actions.AUTH_DID"                   // use to authenticaticate specific did
const AC_SETUP_CHECK = "account.actions.DID_SETUP_CHECK"         // use to check if theres an existing did setup
const AC_SETUP_DID = "account.actions.DID_SETUP"                 // use to setup did
const AC_VALIDATE_TOKEN = PACKAGE_ID + ".actions.VALIDATE_TOKEN" // use to validate token
const AC_AUTH_SYSTEM = PACKAGE_ID + ".actions.AUTH_SYSTEM"       // authenticate via system

const PACKAGE_ID = "ela.account"
const HOME_DIR = "/home/elabox"
const DID_DATA_DIR = HOME_DIR + "/data/" + PACKAGE_ID
const DID_HASH_PATH = DID_DATA_DIR + "/did.dat"
const KEYSTORE_PATH = "/home/elabox/documents/ela.mainchain/keystore.dat"

var ACCOUNT_FILE = path.GetSystemAppDirData(PACKAGE_ID) + "/ac.dat"

var Controller *app.Controller

// error codes
const CREDENTIALS_INVALID = 800
