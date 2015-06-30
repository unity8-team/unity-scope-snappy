package previews

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages"
	"launchpad.net/unity-scope-snappy/store/progress"
	"launchpad.net/unity-scope-snappy/webdm"
)

// NewPreview is a factory for getting the correct preview for a given package.
//
// Parameters:
// snap: Snap to be represented by the preview.
// metadata: Metadata to be used for informing the preview creation.
func NewPreview(snap webdm.Package, metadata *scopes.ActionMetadata) (interfaces.PreviewGenerator, error) {
	// Temporary hack to provide a manual refresh while support for progrss is
	// being added.
	progressHack := &progress.Hack{}
	if retrieveProgressHack(metadata, progressHack) {
		// If an operation is still ongoing, show progress
		if snap.Status != progressHack.DesiredStatus {
			switch progressHack.DesiredStatus {
			case webdm.StatusInstalled:
				return NewInstallingPreview(snap)
			case webdm.StatusNotInstalled:
				return NewUninstallingPreview(snap)
			default:
				return nil, fmt.Errorf("Unexpected desired status: %d", progressHack.DesiredStatus)
			}
		} else {
			if !snap.Installed() && !snap.NotInstalled() {
				return nil, fmt.Errorf("Invalid desired status for progress: %d", progressHack.DesiredStatus)
			}
		}
	}

	return packages.NewPreview(snap)
}

// retrieveProgressHack is used to obtain the ProgressHack struct from
// ActionMetadata.
//
// Parameters:
// metadata: ActionMetadata potentially containing progress hack
// progressHack: Retrieved progress.Hack (if any)
//
// Returns:
// - Whether or not a progress.Hack was retrieved.
func retrieveProgressHack(metadata *scopes.ActionMetadata, progressHack *progress.Hack) bool {
	err := metadata.ScopeData(progressHack)
	return (err == nil) && (progressHack.DesiredStatus != webdm.StatusUndefined)
}
