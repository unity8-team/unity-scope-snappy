package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// WidgetReceiver is an interface to be implemented by any struct that supports
// the type of preview widget interface used by this scope.
type WidgetReceiver interface {
	PushWidgets(widgets ...scopes.PreviewWidget) error
}

// PreviewGenerator is an interface to be implemented by any struct that wishes
// to provide previews for use in this scope.
type PreviewGenerator interface {
	Generate(receiver WidgetReceiver) error
}

// NewPreview is a factory for getting the correct preview for a given package.
//
// Parameters:
// snap: Snap to be represented by the preview.
// metadata: Metadata to be used for informing the preview creation.
func NewPreview(snap webdm.Package, metadata *scopes.ActionMetadata) (PreviewGenerator, error) {
	// Temporary hack to provide a manual refresh while support for progrss is
	// being added.
	progressHack := &ProgressHack{}
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

	if snap.Installed() {
		return NewInstalledPreview(snap)
	} else {
		return NewStorePreview(snap)
	}
}

// retrieveProgressHack is used to obtain the ProgressHack struct from
// ActionMetadata.
//
// Parameters:
// metadata: ActionMetadata potentially containing progress hack
// progressHack: Retrieved ProgresHack (if any)
//
// Returns:
// - Whether or not a ProgressHack was retrieved.
func retrieveProgressHack(metadata *scopes.ActionMetadata, progressHack *ProgressHack) bool {
	err := metadata.ScopeData(progressHack)
	return (err == nil) && (progressHack.DesiredStatus != webdm.StatusUndefined)
}
