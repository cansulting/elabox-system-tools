package main

type Account struct {
	Address string `json:"address"`
}

func RetrievePrimaryAccountDetails() map[string]interface{} {
	result := make(map[string]interface{})
	wallet, _ := LoadWalletAddr()
	result["wallet"] = wallet
	return result
}
