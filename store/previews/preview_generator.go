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

package previews

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages"
	"launchpad.net/unity-scope-snappy/webdm"
)

// NewPreview is a factory for getting the correct preview for a given package.
//
// Parameters:
// snap: Snap to be represented by the preview.
// metadata: Metadata to be used for informing the preview creation.
func NewPreview(snap webdm.Package, metadata *scopes.ActionMetadata) (interfaces.PreviewGenerator, error) {
	var operationMetadata operation.Metadata

	// This may fail, but the zero-value of OperationMetadata is fine
	metadata.ScopeData(&operationMetadata)

	// If an uninstall was requested, per store design we need to confirm the
	// request.
	if operationMetadata.UninstallRequested {
		return NewConfirmUninstallPreview(snap), nil
	}

	return packages.NewPreview(snap, operationMetadata)
}
