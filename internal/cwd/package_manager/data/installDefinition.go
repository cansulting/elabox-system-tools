package data

type InstallDef struct {
	Id   string `json:"id"`   // definition id, usually the package id
	Url  string `json:"url"`  // where the data will be downloaded
	Icon string `json:"icon"` // icon that represents the install def
	Name string `json:"name"`
}

func (inst InstallDef) ToPackageInfo() PackageInfo {
	pkgi := PackageInfo{
		Id:   inst.Id,
		Name: inst.Name,
		Icon: inst.Icon,
	}
	return pkgi
}

func (inst InstallDef) FromMap(data map[string]interface{}) InstallDef {
	inst.Id = data["id"].(string)
	inst.Icon = data["icon"].(string)
	inst.Name = data["name"].(string)
	inst.Url = data["url"].(string)
	return inst
}
