package data

// type that represents a package from store server listing
type PackageListingCache struct {
	Id    string `json:"id"`   // Package ID
	Name  string `json:"name"` // Package name
	Icon  string `json:"icon"` // Package icon url
	Build int    `json:"build"`
	// extra details
	Version        string   `json:"version"`
	Description    string   `json:"description"`
	Updates        string   `json:"updates"`
	ProjectRepo    string   `json:"projectRepo"`
	ProjectWebsite string   `json:"projectWebsite"`
	Category       string   `json:"category"`
	MinRuntime     string   `json:"minRuntime"`             // minimum runtime required to run this package, val <= 0 means it can run on any runtime
	Beta           bool     `json:"beta"`                   // true if this is a beta package
	BetaUsers      []string `json:"betaUsers,omitempty"`    //list of beta users
	Dependencies   []string `json:"dependencies,omitempty"` // list of package ids where this package is dependent to
}
