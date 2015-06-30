package packages

import (
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates"
	"launchpad.net/unity-scope-snappy/webdm"
)

// Preview is a PreviewGenerator representing a given package.
type Preview struct {
	template templates.Template
}

// NewPreview creates a new Preview for representing a given package.
//
// Parameters:
// snap: Package to be represented by the preview.
func NewPreview(snap webdm.Package) (*Preview, error) {
	preview := new(Preview)

	var err error
	if snap.Installed() {
		preview.template, err = templates.NewInstalledTemplate(snap)
	} else {
		preview.template, err = templates.NewStoreTemplate(snap)
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
func (preview Preview) Generate(receiver interfaces.WidgetReceiver) error {
	receiver.PushWidgets(preview.template.HeaderWidget())
	receiver.PushWidgets(preview.template.ActionsWidget())
	receiver.PushWidgets(preview.template.InfoWidget())
	receiver.PushWidgets(preview.template.UpdatesWidget())

	return nil
}
