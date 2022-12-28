package main

import (
	"os"
	"os/exec"
	"strings"
)

const SHADOW_FILE = "/etc/shadow"

// check authorization for specific DID, return true if DID is authorized
// func AuthenticateDid(did string) bool {
// 	deviceSerial := system.GetDeviceInfo().Serial
// 	hash := sha256.Sum256([]byte(did + deviceSerial))
// 	// step: load the currently saved did hash
// 	savedHash, err := os.ReadFile(DID_HASH_PATH)
// 	if err != nil {
// 		return false
// 	}
// 	if string(hash[:]) != string(savedHash) {
// 		return false
// 	}
// 	return true
// }

// use to authenticate specific password
func AuthenticateSystemAccount(username string, password string) (error, bool) {
	contents, err := os.ReadFile(SHADOW_FILE)
	if err != nil {
		return err, false
	}
	hashContent := Grep(username, string(contents))
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
