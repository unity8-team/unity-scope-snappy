package templates

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

// Template is an interface to be implemented by structs which
// are representing a given package for a Unity scope preview.
type Template interface {
	// HeaderWidget generates a widget for the preview header section.
	HeaderWidget() scopes.PreviewWidget

	// ActionsWidget generates a widget for the preview actions section.
	ActionsWidget() scopes.PreviewWidget

	// InfoWidget generates a widget for the preview info section.
	InfoWidget() scopes.PreviewWidget

	// UpdatesWidget generates a widget for the preview updates section.
	UpdatesWidget() scopes.PreviewWidget
}
