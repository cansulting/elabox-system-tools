package data

type PackageConfigError struct {
	propertyError string
}

func (e *PackageConfigError) Error() string {
	return "PackageConfig: " + e.propertyError
}
