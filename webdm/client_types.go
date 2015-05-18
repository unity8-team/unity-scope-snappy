package webdm

// Package contains information about a given package available from the store
// or already installed.
type Package struct {
	Id           string
	Name         string
	Origin       string
	Version      string
	Vendor       string
	Description  string
	IconUrl      string
	Installed    bool
	DownloadSize int
	Type         string
}
