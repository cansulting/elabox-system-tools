package main

import (
	"crypto/sha256"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/utils"
)

const SHADOW_FILE = "/etc/shadow"

// check authorization for specific DID, return true if DID is authorized
func AuthenticateDid(did string) bool {
	deviceSerial := system.GetDeviceInfo().Serial
	hash := sha256.Sum256([]byte(did + deviceSerial))
	// step: load the currently saved did hash
	savedHash, err := os.ReadFile(data.DID_HASH_PATH)
	if err != nil {
		return false
	}
	if string(hash[:]) != string(savedHash) {
		return false
	}
	return true
}

func IsDidSetup() bool {
	if _, err := os.Stat(data.DID_HASH_PATH); err != nil {
		return false
	}
	return true
}

// use to authenticate specific password
func AuthenticateSystemAccount(username string, password string) (error, bool) {
	contents, err := os.ReadFile(SHADOW_FILE)
	if err != nil {
		return err, false
	}
	hashContent := utils.Grep(username, string(contents))
	// unable to find specific user account
	if hashContent == "" {
		return nil, false
	}
	creds := strings.Split(hashContent, "$")
	salt := creds[2]
	savedHash := strings.Split(hashContent, ":")[1]
	encryptType := creds[1]
	// generate hash
	cmd := exec.Command(
		"/usr/bin/openssl",
		"passwd", "-"+encryptType,
		"-salt", salt,
		password,
	)
	hash, err := cmd.CombinedOutput()
	if err != nil {
		return err, false
	}
	// password is correct
	strHash := string(hash)
	strHash = strings.TrimRight(strHash, "\n")
	if savedHash == strHash {
		return nil, true
	}
	return nil, false
}

// set the current device did
func SetDeviceDid(presentation map[string]interface{}) error {
	// step: validate presentation
	if presentation["holder"] == nil {
		return errors.New("no holder provider in presentation")
	}
	// step: create hash
	did := presentation["holder"].(string)
	deviceSerial := system.GetDeviceInfo().Serial
	hash := sha256.Sum256([]byte(did + deviceSerial))
	// step: save to file
	if err := os.MkdirAll(data.DID_DATA_DIR, perm.PUBLIC_WRITE); err != nil {
		return err
	}
	if err := os.WriteFile(data.DID_HASH_PATH, hash[:], perm.PUBLIC_VIEW); err != nil {
		return err
	}
	return nil
}
