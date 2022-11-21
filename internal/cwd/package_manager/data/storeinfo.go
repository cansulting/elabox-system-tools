package data

type StorePreview struct {
	Id          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"desc"`
	CID         string                    `json:"storecid"`
	Packages    map[string]PackagePreview `json:"packages,omitempty"`
}
