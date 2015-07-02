package fakes

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

// FakeWidgetReceiver is a fake implementation of the WidgetReceiver interface,
// for use within tests.
type FakeWidgetReceiver struct {
	Widgets []scopes.PreviewWidget
}

func (receiver *FakeWidgetReceiver) PushWidgets(widgets ...scopes.PreviewWidget) error {
	if receiver.Widgets == nil {
		receiver.Widgets = widgets
	} else {
		receiver.Widgets = append(receiver.Widgets, widgets...)
	}

	return nil
}
