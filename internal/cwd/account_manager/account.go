package main

import (
	"encoding/json"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/perm"
)

var currentAccount *Account = nil

type Account struct {
	Status      string            `json:"status"`
	Token       string            `json:"token"`
	DisplayName string            `json:"displayName"`
	Did         string            `json:"did"`
	Wallets     map[string]string `json:"wallets,omitempty"`
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

func GetAccount() *Account {
	if err := initAccount(); err != nil {
		return nil
	}
	return currentAccount
}

func initAccount() error {
	if currentAccount != nil {
		return nil
	}
	if _, err := os.Stat(ACCOUNT_FILE); err != nil {
		currentAccount = createEmptyAccount()
		return nil
	}
	contents, err := os.ReadFile(ACCOUNT_FILE)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, currentAccount)
}

func SaveAccount() error {
	contents, err := json.Marshal(currentAccount)
	if err != nil {
		return err
	}
	return os.WriteFile(ACCOUNT_FILE, contents, perm.PRIVATE)
}

func UpdateWalletAddress(id string, address string) error {
	if err := initAccount(); err != nil {
		return err
	}
	currentAccount.Wallets[id] = address
	return nil
}

// set the current device did
func UpdateDid(did string) error {
	if err := initAccount(); err != nil {
		return err
	}
	currentAccount.Did = did
	return nil
}

// check authorization for specific DID, return true if DID is authorized
func AuthenticateDid(did string) (error, bool) {
	if err := initAccount(); err != nil {
		return err, false
	}
	return nil, currentAccount.Did == did
}
