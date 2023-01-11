package main

import (
	"testing"

	"github.com/cansulting/elabox-system-tools/foundation/account"
)

func Test_Decryption(t *testing.T) {
	privKey, err := GenerateKeyPair(2048)
	if err != nil {
		t.Error(err)
		return
	}
	privStr := PrivateKeyToBytes(privKey)
	pubStr, err := PublicKeyToBytes(&privKey.PublicKey)
	if err != nil {
		t.Error(err)
		return
	}

	println(string(privStr), " ", string(pubStr))
	priv2, err := BytesToPrivateKey(privStr)
	if err != nil {
		t.Error(err)
		return
	}
	pub2, err := PublicKeyToBytes(&priv2.PublicKey)
	if err != nil {
		t.Error(err)
		return
	}
	if string(pub2) == string(pubStr) {
		t.Error("Should not be used")
	}
}

// test creating and validating token
func Test_Token(t *testing.T) {
	token, err := CreateToken("sameuser", "")
	if err != nil {
		t.Error(err)
		return
	}
	pub, err := GetPublicKey()
	if err != nil {
		t.Error("failed retrieving public key")
		return
	}
	if err := account.SetPublicKeyStr(pub); err != nil {
		t.Error(err)
		return
	}
	valid, err := account.ValidateToken(token)
	if err != nil {
		t.Error("error found while validating token", err)
		return
	}
	if !valid {
		t.Error("failed to validate token")
		return
	}
}

func Test_ResolveToken(t *testing.T) {
	err := account.ResolvePublicKey()
	if err != nil {
		t.Error(err)
	}
}
