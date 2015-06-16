package scope

import (
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
func NewPreview(snap webdm.Package) (PreviewGenerator, error) {
	return NewPackagePreview(snap)
}
