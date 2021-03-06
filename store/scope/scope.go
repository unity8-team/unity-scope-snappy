/* Copyright (C) 2015-2016 Canonical Ltd.
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

package scope

import (
	"fmt"
	"log"

	"github.com/snapcore/snapd/client"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/store/previews"
)

// template for the grid layout of the search results.
const layout = `{
	"schema-version": 1,
	"template": {
		"category-layout": "grid",
		"card-size": "small"
	},
	"components": {
		"title": "title",
		"subtitle": "subtitle",
        "attributes": { "field": "attributes", "max-count": 4 },
		"art" : {
			"field": "art",
            "aspect-ratio": 1.13,
            "fallback": "image://theme/placeholder-app-icon"
		}
    }
}`

// Scope is the struct representing the scope itself.
type Scope struct {
	webdmClient packages.WebdmManager
	dbusClient  *packages.DbusManagerClient
}

// New creates a new Scope using a specific WebDM API URL.
//
// Parameters:
//
// Returns:
// - Pointer to scope (nil if error).
// - Error (nil if none).
func New() (*Scope, error) {
	scope := new(Scope)
	var err error
	scope.webdmClient, err = packages.NewSnapdClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to create WebDM client: %s", err)
	}

	scope.dbusClient = packages.NewDbusManagerClient()
	err = scope.dbusClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to dbus session bus: %s", err)
	}

	return scope, nil
}

func (scope Scope) SetScopeBase(base *scopes.ScopeBase) {
	// Do nothing
}

func (scope Scope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	installedApps := scope.webdmClient.GetInstalledPackages()
	available, err := scope.webdmClient.GetStorePackages(query.QueryString())
	if err != nil {
		return scopeError("unity-scope-snappy: Unable to get package list: %s", err)
	}

	var category *scopes.Category
	category = reply.RegisterCategory("store_packages", "Store Packages", "", layout)

	for _, thisPackage := range available {
		_, ok := installedApps[thisPackage.Name]
		result := packageResult(category, thisPackage, ok)

		if reply.Push(result) != nil {
			// If the push fails, the query was cancelled. No need to continue.
			return nil
		}
	}

	return nil
}

func (scope Scope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	var snapName string
	err := result.Get("name", &snapName)
	if err != nil {
		return scopeError(`unity-scope-snappy: Unable to retrieve ID for package "%s": %s`, result.Title(), err)
	}

	// Need to query the API to make sure we have an up-to-date status,
	// otherwise we can't refresh the state of the buttons after an install or
	// uninstall action.
	snap, err := scope.webdmClient.Query(snapName)

	if err != nil {
		return scopeError(`unity-scope-snappy: Unable to query API for package "%s": %s`, result.Title(), err)
	}

	preview, err := previews.NewPreview(*snap, result, metadata)
	if err != nil {
		return scopeError(`unity-scope-snappy: Unable to create preview for package "%s": %s`, result.Title(), err)
	}

	err = preview.Generate(reply)
	if err != nil {
		return scopeError(`unity-scope-snappy: Unable to generate preview for package "%s": %s`, result.Title(), err)
	}

	return nil
}

func (scope *Scope) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
	// Obtain the ID for the specific package
	var snapId string
	err := result.Get("name", &snapId)
	if err != nil {
		return nil, scopeError(`unity-scope-snappy: Unable to retrieve ID for package "%s": %s`, result.Title(), err)
	}

	// Get the action runner corresponding to this action
	actionRunner, err := actions.NewRunner(actions.ActionId(actionId))
	if err != nil {
		return nil, scopeError(`unity-scope-snappy: Unable to handle action "%s": %s`, actionId, err)
	}

	response, err := actionRunner.Run(scope.dbusClient, snapId)
	if err != nil {
		err = scopeError(`unity-scope-snappy: Error handling action "%s": %s`, actionId, err)
	}

	return response, err
}

// packageResult is used to create a scopes.CategorisedResult from a
// client.Snap.
//
// Parameters:
// category: Category in which the result will be created.
// snap: client.Snap representing snap.
//
// Returns:
// - Pointer to scopes.CategorisedResult
func packageResult(category *scopes.Category, snap client.Snap, installed bool) *scopes.CategorisedResult {
	result := scopes.NewCategorisedResult(category)

	// NOTE: Title really needs to be title, not name, but snapd doesn't expose
	result.SetTitle(snap.Name)
	result.SetArt(snap.Icon)
	result.SetURI("snappy:" + snap.Name)
	result.Set("subtitle", snap.Developer)
	result.Set("name", snap.Name)
	result.Set("id", snap.ID)
	result.Set("installed", installed)
	var price string
	if installed == true {
		price = "✔ INSTALLED"
	} else {
		price = "FREE"
	}
	result.Set("price_area", price)
	// This is a bit of a mess at the moment, need a better way to do this
	attributes := make([]map[string]string, 0)
	emptyValue := make(map[string]string, 0)
	emptyValue["value"] = ""
	priceValue := make(map[string]string, 0)
	priceValue["value"] = price
	attributes = append(attributes, priceValue)
	attributes = append(attributes, emptyValue)
	attributes = append(attributes, emptyValue)
	attributes = append(attributes, emptyValue)
	result.Set("attributes", attributes)
	return result
}

// scopeError prints an error to stderr as well as returning an actual error.
// This is used because the errors returned from the scope functions don't seem
// to be handled, logged, or otherwise displayed to the user in any way.
func scopeError(format string, a ...interface{}) error {
	log.Printf(format, a...)
	return fmt.Errorf(format, a...)
}
