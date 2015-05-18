package webdm

type Status bool

// Package contains information about a given package available from the store
// or already installed.
type Package struct {
	Id           string
	Name         string
	Origin       string
	Version      string
	Vendor       string
	Description  string
	IconUrl      string `json:"icon"`
	Installed    Status `json:"status"`
	DownloadSize int    `json:"download_size"`
	Type         string
}
