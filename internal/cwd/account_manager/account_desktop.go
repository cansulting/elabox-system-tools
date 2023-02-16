//go:build !linux
// +build !linux

package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func changeSystemPassword(pass string) error {
	// dont change system password on desktop pc
	return nil
}

func AuthenticateSystemPassword(username string, password string) (error, bool) {
	acc, err := GetAccount(username)
	if err != nil {
		return err, false
	}
	if password == "" {
		return errors.New("password was empty"), false
	}
	sum := md5.Sum([]byte(password))
	passhash := hex.EncodeToString(sum[:])
	if acc.PassHash == passhash {
		return nil, true
	}
	return nil, false
}
