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
