package data

import (
	"path"
	"strconv"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
)

type PackageInfo struct {
	Id               string           `json:"id"`   // Package ID
	Name             string           `json:"name"` // Package name
	Icon             string           `json:"icon"` // Package icon
	CurrentBuild     int              `json:"currentBuild"`
	LatestBuild      int              `json:"latestBuild"`
	Status           global.AppStatus `json:"status"`
	Progress         float32          `json:"progress"`
	Notifications    int              `json:"notifications"`
	Description      string           `json:"description,omitempty"`
	Updates          string           `json:"updates,omitempty"`
	Version          string           `json:"version,omitempty"`
	LaunchUrl        string           `json:"launchUrl,omitempty"`
	Category         string           `json:"category,omitempty"`
	IsService        bool             `json:"isService"`
	LatestMinRuntime string           `json:"latestMinRuntime,omitempty"` // the minimum runtime required to install this package
	Dependencies     []string         `json:"dependencies,omitempty"`     // list of package ids where this package is dependent to
	IsDependency     bool             `json:"isDependent,omitempty"`
	Enabled          bool             `json:"enabled,omitempty"`
}

func NewPackageInfo() PackageInfo {
	return PackageInfo{}
}

// add informations
func (instance *PackageInfo) AddInfo(installed *data.PackageConfig, storeCacheItem *PackagePreview, detailed bool) {
	if installed != nil {
		instance.Id = installed.PackageId
		instance.Name = installed.Name
		instance.LatestBuild = int(installed.Build)
		instance.CurrentBuild = int(installed.Build)
		instance.IsService = installed.HasServices()
		if installed.Icons != nil {
			instance.Icon = installed.GetIconUrl() + "/" + path.Base(installed.Icons.Medium)
		}
		if !installed.HasServices() ||
			installed.ActivityGroup.CustomLink != "" ||
			installed.ActivityGroup.CustomPort != 0 {
			instance.LaunchUrl = "/" + installed.PackageId
			if installed.ActivityGroup.CustomLink != "" {
				instance.LaunchUrl = installed.ActivityGroup.CustomLink
			} else {
				if installed.ActivityGroup.CustomPort != 0 {
					instance.LaunchUrl = ":" + strconv.Itoa(installed.ActivityGroup.CustomPort)
				}
			}
		}
	}
	if storeCacheItem != nil {
		instance.Id = storeCacheItem.Id
		instance.Name = storeCacheItem.Name
		instance.LatestBuild = int(storeCacheItem.Release.Production.Build.Number)
		instance.Icon = storeCacheItem.IconCID
		if detailed {
			// if instance.Details == nil {
			// 	instance.Details = &PackageDetails{}
			// }
			instance.LatestMinRuntime = storeCacheItem.Release.Production.Build.MinRuntime
			//instance.Category = storeCacheItem.Category
			instance.Description = storeCacheItem.Description
			instance.Updates = storeCacheItem.Release.Production.Description
			instance.Version = storeCacheItem.Release.Production.Version
			instance.Dependencies = storeCacheItem.Release.Production.Build.Dependencies
			instance.IsDependency = false
			// if loaded, _ := storeCacheItem.LoadDetails(); loaded {

			// }
		}
	}

	if instance.CurrentBuild == 0 {
		instance.Status = "uninstalled"
	} else {
		instance.Status = "installed"
	}
}
