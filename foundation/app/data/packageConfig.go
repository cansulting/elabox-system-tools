package data

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type PackageConfig struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	PackageId        string            `json:"packageId"`      // company.package
	Build            int16             `json:"build"`          // major.minor.patch
	Version          string            `json:"version"`        // current version
	Permissions      []string          `json:"permissions"`    // declared actions to be called inside the app
	Services         map[string]string `json:"services"`       // if app has a service. this contains definition of commands available to service
	Activities       []string          `json:"activities"`     // if app has activity. this contains definition of actions that will triggerr activity
	BroacastListener []string          `json:"actionListener"` // defined actions which action listener will listen to
	InstallLocation  string            `json:"location"`       // which location the package will be installed
	Source           string            `json:"-"`
}

func DefaultPackage() *PackageConfig {
	return &PackageConfig{InstallLocation: "external"}
}

func (c *PackageConfig) LoadFromSrc(src string) error {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return &PackageConfigError{propertyError: "Error loading package. " + err.Error()}

	}
	c.Source = src
	return c.LoadFromBytes(bytes)
}

func (c *PackageConfig) LoadFromReader(reader io.Reader) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return &PackageConfigError{propertyError: err.Error()}
	}
	err = c.LoadFromBytes(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *PackageConfig) LoadFromBytes(bytes []byte) error {
	return json.Unmarshal(bytes, &c)
}

func (c *PackageConfig) GetError() error {
	if c.PackageId == "" {
		return &PackageConfigError{propertyError: "Package Id shouldnt be null"}
	}
	return nil
}
