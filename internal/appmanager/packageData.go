package main

type PackageData struct {
	AppId   string                 `json:"appId"`
	Actions []ActionDefinitionData `json:"actions"`
}
