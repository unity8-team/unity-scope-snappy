package scope

import (
	"launchpad.net/unity-scope-snappy/webdm"
)

// PackagePreview is a PreviewGenerator representing a given package.
type PackagePreview struct {
	template PackagePreviewTemplate
}

// NewPackagePreview creates a new PackagePreview for representing a given
// package.
//
// Parameters:
// snap: Package to be represented by the preview.
func NewPackagePreview(snap webdm.Package) (*PackagePreview, error) {
	preview := new(PackagePreview)

	var err error
	if snap.Installed() {
		preview.template, err = NewInstalledPackagePreviewTemplate(snap)
	} else {
		preview.template, err = NewStorePackagePreviewTemplate(snap)
	}

	return preview, err
}

// Generate pushes the template's preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview PackagePreview) Generate(receiver WidgetReceiver) error {
	receiver.PushWidgets(preview.template.headerWidget())
	receiver.PushWidgets(preview.template.actionsWidget())
	receiver.PushWidgets(preview.template.infoWidget())
	receiver.PushWidgets(preview.template.updatesWidget())

	return nil
}
