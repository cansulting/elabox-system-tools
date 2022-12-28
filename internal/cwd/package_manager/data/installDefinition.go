package data

type InstallDef struct {
	Id   string `json:"id"`
	Url  string `json:"url"`
	Icon string `json:"icon"`
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
