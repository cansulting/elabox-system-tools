package data

type ReleaseType uint

const (
	Development ReleaseType = 0
	Beta        ReleaseType = 1
	Production  ReleaseType = 2
)

type BuildInfo struct {
	Id           string   `json:"id"`
	Number       uint     `json:"number"`
	IpfsCID      string   `json:"ipfsCID"`
	MinRuntime   string   `json:"minRuntime"`
	Dependencies []string `json:"dependencies"`
}

type ReleaseUnit struct {
	Description string      `json:"desc"`
	Build       BuildInfo   `json:"build"`
	ReleaseType ReleaseType `json:"type"`
	Version     string      `json:"version"`
	DateEpoch   int         `json:"dateEpoch"`
	Users       []string    `json:"testers,omitempty"`
}

type ReleaseInfo struct {
	Alpha      ReleaseUnit `json:"alpha,omitempty"`
	Beta       ReleaseUnit `json:"beta,omitempty"`
	Production ReleaseUnit `json:"prod,omitempty"`
}

type PackagePreview struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	IconCID     string      `json:"iconcid"`
	Release     ReleaseInfo `json:"release,omitempty"`
}
