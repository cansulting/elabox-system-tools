package main

type InstallError struct {
	errorString string
}

func (e *InstallError) Error() string {
	return e.errorString
}
