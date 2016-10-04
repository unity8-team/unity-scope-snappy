/* Copyright (C) 2015 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

package fakes

import (
	"launchpad.net/go-unityscopes/v2"
	"reflect"
	"testing"
)

// Test pushing a single widget onto FakeWidgetReceiver
func TestFakeWidgetReceiver_singleWidget(t *testing.T) {
	receiver := FakeWidgetReceiver{}

	widget := scopes.NewPreviewWidget("header", "header")
	receiver.PushWidgets(widget)

	if len(receiver.Widgets) != 1 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Widget list contained %d items, expected 1", len(receiver.Widgets))
	}

	if !reflect.DeepEqual(receiver.Widgets[0], widget) {
		t.Error("Expected widget list to contain the header widget we gave it")
	}
}

// Test pushing multiple widgets onto FakeWidgetReceiver
func TestFakeWidgetReceiver_multipleWidgets(t *testing.T) {
	receiver := FakeWidgetReceiver{}

	widget1 := scopes.NewPreviewWidget("header", "header")
	widget2 := scopes.NewPreviewWidget("text", "text")
	receiver.PushWidgets(widget1, widget2)

	if len(receiver.Widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Widget list contained %d items, expected 2", len(receiver.Widgets))
	}

	// Order is enforced, so these should be predictable
	if !reflect.DeepEqual(receiver.Widgets[0], widget1) {
		t.Error("Expected widget list to contain the header widget we gave it")
	}

	if !reflect.DeepEqual(receiver.Widgets[1], widget2) {
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

	if len(receiver.Widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Widget list contained %d items, expected 2", len(receiver.Widgets))
	}

	// Order is enforced, so these should be predictable
	if !reflect.DeepEqual(receiver.Widgets[0], widget1) {
		t.Error("Expected widget list to contain the header widget we gave it")
	}

	if !reflect.DeepEqual(receiver.Widgets[1], widget2) {
		t.Error("Expected widget list to contain the text widget we gave it")
	}
}
