package appman

type AppmanError struct {
	errorStr string
}

func (p *AppmanError) Error() string {
	return "PackageManager: " + p.errorStr
}
