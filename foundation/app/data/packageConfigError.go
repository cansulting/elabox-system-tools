package data

type PackageConfigError struct {
	propertyError string
}

func (e *PackageConfigError) Error() string {
	return "Invalid package " + e.propertyError
}
