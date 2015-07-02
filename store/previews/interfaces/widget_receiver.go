package interfaces

import "launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"

// WidgetReceiver is an interface to be implemented by any struct that supports
// the type of preview widget interface used by this scope.
type WidgetReceiver interface {
	PushWidgets(widgets ...scopes.PreviewWidget) error
}
