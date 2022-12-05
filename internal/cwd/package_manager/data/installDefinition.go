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
