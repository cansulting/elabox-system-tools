package data

import "github.com/cansulting/elabox-system-tools/foundation/app"

// actions
const AC_AUTH_DID = "account.actions.AUTH_DID"           // use to authenticaticate specific did
const AC_SETUP_CHECK = "account.actions.DID_SETUP_CHECK" // use to check if theres an existing did setup
const AC_SETUP_DID = "account.actions.DID_SETUP"         // use to setup did

const PACKAGE_ID = "ela.account"
const HOME_DIR = "/home/elabox"
const DID_DATA_DIR = HOME_DIR + "/data/" + PACKAGE_ID
const DID_HASH_PATH = DID_DATA_DIR + "/did.dat"
const KEYSTORE_PATH = "/home/elabox/documents/ela.mainchain/keystore.dat"

var Controller *app.Controller
