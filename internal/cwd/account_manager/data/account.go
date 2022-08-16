package data

type Account struct {
	Address             string `json:"address"`
	Username            string
	PermissionProfileId string          // empty if this is custom profile
	Permissions         GroupPermission // empty if theres a ref permision profile
}
