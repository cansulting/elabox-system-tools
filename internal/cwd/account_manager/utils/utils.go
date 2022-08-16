package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/data"
)

type Address struct {
	Address     string `json:"Address"`
	ProgramHash string `json:"ProgramHash"`
}

type KeyStore struct {
	Account []Address `json:"Account"`
}

func Grep(keyword string, src string) string {
	splits := strings.Split(src, "\n")
	keywordbt := []byte(keyword)
	for _, line := range splits {
		if len(line) >= len(keywordbt) {
			found := true
			for i, keywordC := range keywordbt {
				if keywordC != line[i] {
					found = false
					break
				}
			}
			if found {
				return line
			}
		}
	}
	// nothing was found
	return ""
}

// use to retrieve wallet address from keystore path
func LoadWalletAddr() (string, error) {
	dat, err := os.ReadFile(data.KEYSTORE_PATH)
	if err != nil {
		return "", err
	}

	var keyStore KeyStore
	_ = json.Unmarshal(dat, &keyStore)
	return keyStore.Account[0].Address, nil
}
