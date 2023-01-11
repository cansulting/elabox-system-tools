package account

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/dgrijalva/jwt-go"
)

var publicKey *rsa.PublicKey

func SetPublicKeyStr(val []byte) error {
	pub, err := jwt.ParseRSAPublicKeyFromPEM(val)
	publicKey = pub
	return err
}

func SetPublicKey(pub *rsa.PublicKey) {
	publicKey = pub
}

func ValidateTokenFromAction(ac data.Action) (bool, error) {
	params, err := ac.DataToMap()
	if err != nil {
		return false, err
	}
	if params["token"] == nil {
		return false, errors.New("no token was provided")
	}
	return ValidateToken(params["token"].(string))
}

func ValidateToken(tokenStr string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(_token *jwt.Token) (interface{}, error) {
		if _, ok := _token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

// resolve public key by requesting via account rpc
func ResolvePublicKey() error {
	if publicKey != nil {
		return nil
	}
	handler, err := rpc.GetInstance()
	if err != nil {
		return errors.New("failed to initialize RPC")
	}
	res, err := handler.CallRPC(
		constants.ACCOUNT_SYS_ID,
		data.NewActionById(constants.ACTION_RETRIEVE_PUBKEY))
	if err != nil {
		return errors.New("failed to execute rpc")
	}
	resmsg, err := res.ToSimpleResponse()
	if err != nil {
		return errors.New("failed to execute rpc. " + err.Error())
	}
	decoded, err := base64.StdEncoding.DecodeString(resmsg.Message)
	if err != nil {
		return errors.New("failed decoding response msg, " + string(decoded))
	}
	if err := SetPublicKeyStr([]byte(decoded)); err != nil {
		return err
	}
	return nil
}
