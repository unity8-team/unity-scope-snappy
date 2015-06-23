package store

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

type PackagePreviewTemplate interface {
	// headerWidget generates a widget for the preview header.
	headerWidget() scopes.PreviewWidget

	// actionsWidget generates a widget for the preview actions.
	actionsWidget() scopes.PreviewWidget

	// infoWidget generates a widget for the preview info section.
	infoWidget() scopes.PreviewWidget

	// updatesWidget generates a widget for the preview updates section.
	updatesWidget() scopes.PreviewWidget
}
