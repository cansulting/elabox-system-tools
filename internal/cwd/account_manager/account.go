package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/notify"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/foundation/system"
)

var lastActive *Account = nil

type Account struct {
	PassHash    string            `json:"pass"`
	UserName    string            `json:"username"`
	Status      string            `json:"status"`
	Token       string            `json:"token"`
	DisplayName string            `json:"displayName"`
	Did         string            `json:"did"`
	Wallets     map[string]string `json:"wallets,omitempty"`
}

func init() {
	if err := os.MkdirAll(ACCOUNT_LOC, perm.PRIVATE); err != nil {
		logger.GetInstance().Error().Err(err).Stack().Msg("failed creating directory")
	}
}

func (instance Account) ToJson() string {
	content, _ := json.Marshal(instance)
	return string(content)
}

func createEmptyAccount() *Account {
	return &Account{
		Wallets: make(map[string]string),
	}
}

func getAccountPath(username string) string {
	return ACCOUNT_LOC + "/" + username
}

func SetupAccount(username string, pass string, displayName string) (*Account, error) {
	if err := initAccount(); err != nil {
		return nil, err
	}
	if pass == "" || username == "" {
		return nil, errors.New("password or username is required")
	}
	acc := createEmptyAccount()
	hashPass := md5.Sum([]byte(pass))
	acc.UserName = username
	acc.PassHash = hex.EncodeToString(hashPass[:])
	acc.DisplayName = displayName
	lastActive = acc
	if err := saveAccount(acc); err != nil {
		return nil, err
	}
	// STEP: change the system password if not yet configured. this is admin
	if !system.IsConfig() {
		changeSystemPassword(pass)
	}
	return acc, nil
}

func GetAccount(username string) (*Account, error) {
	if err := initAccount(); err != nil {
		return nil, err
	}
	// is this the last active, then return
	if lastActive != nil && lastActive.UserName == username {
		return lastActive, nil
	}
	// not exist? return empty
	if _, err := os.Stat(getAccountPath(username)); err != nil {
		return nil, nil
	}
	contents, err := os.ReadFile(getAccountPath(username))
	if err != nil {
		return nil, err
	}
	acc := &Account{}
	err2 := json.Unmarshal(contents, acc)
	if err2 == nil {
		lastActive = acc
	}
	return acc, err2
}

func initAccount() error {
	return nil
}

func saveAccount(acc *Account) error {
	contents, err := json.Marshal(acc)
	if err != nil {
		return err
	}
	return os.WriteFile(getAccountPath(acc.UserName), contents, perm.PRIVATE)
}

func UpdatePassword(username string, oldPass string, newPass string) error {
	acc, err := GetAccount(username)
	if err != nil {
		return err
	}
	// validate from old password
	hash := md5.Sum([]byte(oldPass))
	if acc.PassHash != hex.EncodeToString(hash[:]) {
		return errors.New("Invalid password")
	}
	// update password
	hash = md5.Sum([]byte(newPass))
	hashStr := hex.EncodeToString(hash[:])
	if hashStr == acc.PassHash {
		return errors.New("new password shouldn't be the same with old password")
	}
	acc.PassHash = hashStr
	err = saveAccount(acc)
	if err == nil {
		if err := notify.System(notify.NotificationData{
			Title:   "System",
			Message: "System password was changed.",
		}); err != nil {
			logger.GetInstance().Err(err).Msg("failed to notify password changed")
		}
	}
	return err
}

// use to update wallet address for specific user
func UpdateWalletAddress(username string, walletId string, address string) error {
	if err := initAccount(); err != nil {
		return err
	}
	acc, err := GetAccount(username)
	if err != nil {
		return err
	}
	if acc.Wallets == nil {
		acc.Wallets = make(map[string]string)
	}
	acc.Wallets[walletId] = address
	return saveAccount(acc)
}

// set the current device did
func UpdateDid(username string, did string) error {
	if err := initAccount(); err != nil {
		return err
	}
	acc, err := GetAccount(username)
	if err != nil {
		return err
	}
	acc.Did = did
	return saveAccount(acc)
}

func UpdateAccount(presentation map[string]interface{}) error {
	if presentation["holder"] == nil {
		return errors.New("no holder provider in presentation")
	}
	did := presentation["holder"].(string)
	if err := UpdateDid(DEFAULT_USERNAME, did); err != nil {
		return errors.New("failed to setup did, " + err.Error())
	}
	// step: update wallet address
	if presentation["esc"] == nil {
		return errors.New("no esc wallet address provided")
	}
	if err := UpdateWalletAddress(DEFAULT_USERNAME, "esc", presentation["esc"].(string)); err != nil {
		return errors.New("failed updating esc wallet, " + err.Error())
	}
	return nil
}

// check authorization for specific DID, return true if DID is authorized
func AuthenticateDid(username string, did string) (error, bool) {
	if err := initAccount(); err != nil {
		return err, false
	}
	acc, err := GetAccount(username)
	if err != nil {
		return err, false
	}
	return nil, acc.Did == did
}
