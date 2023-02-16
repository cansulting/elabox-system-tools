//go:build linux
// +build linux

package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
)

func changeSystemPassword(pass string) error {
	logger.GetInstance().Debug().Msg("updating system password")
	stdinPass := "elabox:" + pass
	cmd := exec.Command("chpasswd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdinPass)
	}()
	_, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New("failed setting up system password. " + err.Error())
	}
	return nil
}

// use to authenticate specific password
func AuthenticateSystemPassword(username string, password string) (error, bool) {
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
