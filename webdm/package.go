package webdm

import (
	"fmt"
)

type Status int

const (
	StatusUndefined Status = iota
	StatusNotInstalled
	StatusInstalled
	StatusInstalling
	StatusUninstalling
)

// Package contains information about a given package available from the store
// or already installed.
type Package struct {
	Id          string
	Name        string
	Origin      string
	Version     string
	Vendor      string
	Description string
	IconUrl     string `json:"icon"`
	Type        string

	Progress    float64
	Status      Status

	// WebDM uses this field to report errors, etc. It's not always filled.
	Message string

	// InstalledSize will be filled if the package is installed, otherwise
	// DownloadSize will be filled.
	InstalledSize int64 `json:"installed_size"`
	DownloadSize  int64 `json:"download_size"`
}

// UnmarshallJSON exists to decode the Status field from JSON to our enum.
func (status *Status) UnmarshalJSON(data []byte) error {
	if status == nil {
		return fmt.Errorf("UnmarshalJSON: Called on nil pointer")
	}

	dataString := string(data)
	switch dataString {
	case `"uninstalled"`:
		*status = StatusNotInstalled
	case `"installed"`:
		*status = StatusInstalled
	case `"installing"`:
		*status = StatusInstalling
	case `"uninstalling"`:
		*status = StatusUninstalling
	default:
		*status = StatusUndefined
		return fmt.Errorf("UnmarshalJSON: Unhandled Status type: %s", dataString)
	}

	return nil
}

// MarshalJSON exists to encode the Status enum to JSON.
func (status Status) MarshalJSON() ([]byte, error) {
	switch status {
	case StatusNotInstalled:
		return []byte(`"uninstalled"`), nil
	case StatusInstalled:
		return []byte(`"installed"`), nil
	case StatusInstalling:
		return []byte(`"installing"`), nil
	case StatusUninstalling:
		return []byte(`"uninstalling"`), nil
	default:
		return nil, fmt.Errorf("MarshalJSON: Unhandled Status type: %d", status)
	}
}

// Installed is used to check whether or not the package is installed.
func (snap Package) Installed() bool {
	return snap.Status == StatusInstalled
}

// Installing is used to check whether or not the package is being installed.
func (snap Package) Installing() bool {
	return snap.Status == StatusInstalling
}

// NotInstalled is used to check whether or not the package not installed.
func (snap Package) NotInstalled() bool {
	return snap.Status == StatusNotInstalled
}

// Uninstalling is used to check whether or not the package is being uninstalled.
func (snap Package) Uninstalling() bool {
	return snap.Status == StatusUninstalling
}
