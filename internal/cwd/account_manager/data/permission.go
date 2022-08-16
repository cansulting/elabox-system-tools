package data

type GroupPermission struct {
	Id             string // which permission profile is this
	SystemSettings bool
	Package        PackagePermission
}

type SystemPermission struct {
	AllowSsh    bool
	AllowUpdate bool
}

type PackagePermission struct {
	AllowUninstall      bool
	AllowServiceRestart bool
	AllowClearData      bool
	AllowedPackages     []string // 0 length if allow all
}
