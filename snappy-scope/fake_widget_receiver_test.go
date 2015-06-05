package main

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"reflect"
	"testing"
)

// FakeWidgetReceiver is a fake implementation of the WidgetReceiver interface,
// for use within tests.
type FakeWidgetReceiver struct {
	widgets []scopes.PreviewWidget
}

func (receiver *FakeWidgetReceiver) PushWidgets(widgets ...scopes.PreviewWidget) error {
	if receiver.widgets == nil {
		receiver.widgets = widgets
	} else {
		receiver.widgets = append(receiver.widgets, widgets...)
	}

	return nil
}

// Test pushing a single widget onto FakeWidgetReceiver
func TestFakeWidgetReceiver_singleWidget(t *testing.T) {
	receiver := FakeWidgetReceiver{}

	widget := scopes.NewPreviewWidget("header", "header")
	receiver.PushWidgets(widget)

	if len(receiver.widgets) != 1 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Widget list contained %d items, expected 1", len(receiver.widgets))
	}

	if !reflect.DeepEqual(receiver.widgets[0], widget) {
		t.Error("Expected widget list to contain the header widget we gave it")
	}
}

// Test pushing multiple widgets onto FakeWidgetReceiver
func TestFakeWidgetReceiver_multipleWidgets(t *testing.T) {
	receiver := FakeWidgetReceiver{}

	widget1 := scopes.NewPreviewWidget("header", "header")
	widget2 := scopes.NewPreviewWidget("text", "text")
	receiver.PushWidgets(widget1, widget2)

	if len(receiver.widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Widget list contained %d items, expected 2", len(receiver.widgets))
	}

	// Order is enforced, so these should be predictable
	if !reflect.DeepEqual(receiver.widgets[0], widget1) {
		t.Error("Expected widget list to contain the header widget we gave it")
	}

	if !reflect.DeepEqual(receiver.widgets[1], widget2) {
		t.Error("Expected widget list to contain the text widget we gave it")
	}
}

// Test multiple pushes onto FakeWidgetReceiver
func TestFakeWidgetReceiver_multiplePushes(t *testing.T) {
	receiver := FakeWidgetReceiver{}

	widget1 := scopes.NewPreviewWidget("header", "header")
	receiver.PushWidgets(widget1)

	widget2 := scopes.NewPreviewWidget("text", "text")
	receiver.PushWidgets(widget2)

	if len(receiver.widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Widget list contained %d items, expected 2", len(receiver.widgets))
	}

	// Order is enforced, so these should be predictable
	if !reflect.DeepEqual(receiver.widgets[0], widget1) {
		t.Error("Expected widget list to contain the header widget we gave it")
	}

	if !reflect.DeepEqual(receiver.widgets[1], widget2) {
		t.Error("Expected widget list to contain the text widget we gave it")
	}
}
