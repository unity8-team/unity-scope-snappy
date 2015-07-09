package previews

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages"
	"launchpad.net/unity-scope-snappy/webdm"
)

// NewPreview is a factory for getting the correct preview for a given package.
//
// Parameters:
// snap: Snap to be represented by the preview.
// metadata: Metadata to be used for informing the preview creation.
func NewPreview(snap webdm.Package, metadata *scopes.ActionMetadata) (interfaces.PreviewGenerator, error) {
	var operationMetadata operation.Metadata

	// This may fail, but the zero-value of OperationMetadata is fine
	metadata.ScopeData(&operationMetadata)

	// If an uninstall was requested, per store design we need to confirm the
	// request.
	if operationMetadata.UninstallRequested {
		return NewConfirmUninstallPreview(snap), nil
	}

	return packages.NewPreview(snap, operationMetadata)
}
