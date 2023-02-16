package main

import (
	"crypto/rsa"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/account"
	"github.com/dgrijalva/jwt-go"
)

const SHADOW_FILE = "/etc/shadow"
const SIGNING_METHOD = "RS256"

var privateKey *rsa.PrivateKey

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

func init() {
	priv, err := GenerateKeyPair(1024)
	if err != nil {
	}
	privateKey = priv
	account.SetPublicKey(&priv.PublicKey)
}

// use to create token for a given user
func CreateToken(username string, origin string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(SIGNING_METHOD))
	token.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		Subject:   username,
	}
	return token.SignedString(privateKey)
}

func GetPublicKey() ([]byte, error) {
	return PublicKeyToBytes(&privateKey.PublicKey)
}
