package main

import (
	"flag"
	"fmt"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"log"
)

type SnappyScope struct {
	webdmClient *webdm.Client
}

func (scope *SnappyScope) SetScopeBase(base *scopes.ScopeBase) {
	// Do nothing
}

const template = `{
	"schema-version": 1,
	"template": {
		"category-layout": "grid",
		"card-size": "small"
	},
	"components": {
		"title": "title",
		"art" : {
			"field": "art"
		},
		"subtitle": "subtitle"
        }
}`

func (scope SnappyScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	createDepartments(query, reply)

	packages, err := scope.getPackageList(query.DepartmentID())
	if err != nil {
		return scopeError("unity-scope-snappy: Unable to get package list: %s", err)
	}

	var category *scopes.Category
	if query.DepartmentID() == "installed" {
		category = reply.RegisterCategory("installed_packages", "Installed Packages", "", template)
	} else {
		category = reply.RegisterCategory("store_packages", "Store Packages", "", template)
	}

	for _, thisPackage := range packages {
		result, err := packageResult(category, thisPackage)
		if err != nil {
			return scopeError(`unity-scope-snappy: Unable to create result for package "%s": %s`, err)
		}

		if reply.Push(result) != nil {
			// If the push fails, the query was cancelled. No need to continue.
			return nil
		}
	}

	return nil
}

func (scope SnappyScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	var snapId string
	err := result.Get("id", &snapId)
	if err != nil {
		return scopeError(`unity-scope-snappy: Unable to retrieve ID for package "%s": %s`, result.Title(), err)
	}

	// Need to query the API to make sure we have an up-to-date status,
	// otherwise we can't refresh the state of the buttons after an install or
	// uninstall action.
	snap, err := scope.webdmClient.Query(snapId)
	if err != nil {
		return scopeError(`unity-scope-snappy: Unable to query API for package "%s": %s`, result.Title(), err)
	}

	reply.PushWidgets(packageHeaderWidget(*snap))
	reply.PushWidgets(packageActionWidget(*snap))

	return nil
}

func (scope *SnappyScope) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
	var snapId string
	err := result.Get("id", &snapId)
	if err != nil {
		return nil, scopeError(`unity-scope-snappy: Unable to retrieve ID for package "%s": %s`, result.Title(), err)
	}

	switch actionId {
	case "install":
		err = scope.webdmClient.Install(snapId)
		if err != nil {
			return nil, scopeError(`unity-scope-snappy: Unable to install package "%s": %s`, result.Title(), err)
		}

	case "uninstall":
		err = scope.webdmClient.Uninstall(snapId)
		if err != nil {
			return nil, scopeError(`unity-scope-snappy: Unable to uninstall package "%s": %s`, result.Title(), err)
		}

	case "open":
		return nil, scopeError(`unity-scope-snappy: Unable to open package "%s": Opening snaps is not yet supported`, result.Title())
	}

	return scopes.NewActivationResponse(scopes.ActivationShowPreview), nil
}

// getPackageList is used to obtain a package list for a specific department.
//
// Parameters:
// department: The department whose packages should be listed.
//
// Returns:
// - List of WebDM Package structs
// - Error (nil if none)
func (scope SnappyScope) getPackageList(department string) ([]webdm.Package, error) {
	var packages []webdm.Package
	var err error

	switch department {
	case "installed":
		packages, err = scope.webdmClient.GetInstalledPackages()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve installed packages: %s", err)
		}

	default:
		packages, err = scope.webdmClient.GetStorePackages()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve store packages: %s", err)
		}
	}

	return packages, nil
}

func main() {
	webdmAddressParameter := flag.String("webdm", webdm.DefaultApiUrl, "WebDM address[:port]")
	flag.Parse()

	scope := new(SnappyScope)
	var err error
	scope.webdmClient, err = webdm.NewClient(*webdmAddressParameter)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to create webdm client: %s", err)
		return
	}

	err = scopes.Run(scope)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to run scope: %s", err)
	}
}

// createDepartments is used to create a set of static departments for the scope.
//
// Parameters:
// query: Query to be executed when the department is selected.
// reply: Reply onto which the departments will be registered
//
// Returns:
// - Error (nil if none)
func createDepartments(query *scopes.CannedQuery, reply *scopes.SearchReply) error {
	rootDepartment, err := scopes.NewDepartment("", query, "All Categories")
	if err != nil {
		return fmt.Errorf(`Unable to create "All Categories" department: %s`, err)
	}

	installedDepartment, err := scopes.NewDepartment("installed", query, "My Snaps")
	if err != nil {
		return fmt.Errorf(`Unable to create "My Snaps" department: %s`, err)
	}

	rootDepartment.SetSubdepartments([]*scopes.Department{installedDepartment})
	reply.RegisterDepartments(rootDepartment)

	return nil
}

// packageResult is used to create a scopes.CategorisedResult from a
// webdm.Package.
//
// Parameters:
// category: Category in which the result will be created.
// snap: webdm.Package representing snap.
//
// Returns:
// - Pointer to scopes.CategorisedResult
// - Error (nil if none)
func packageResult(category *scopes.Category, snap webdm.Package) (*scopes.CategorisedResult, error) {
	result := scopes.NewCategorisedResult(category)

	result.SetTitle(snap.Name)
	result.Set("subtitle", snap.Description)
	result.SetURI("snappy:" + snap.Id)
	result.SetArt(snap.IconUrl)

	// The preview needs access to the ID
	err := result.Set("id", snap.Id)
	if err != nil {
		return nil, fmt.Errorf(`Unable to set ID for package "%s": %s`, snap.Name, err)
	}

	return result, nil
}

// packageHeaderWidget is used to create a header widget to preview a single snap.
//
// Parameters:
// snap: webdm.Package representing the snap.
//
// Returns:
// - Header preview widget for the given snap.
func packageHeaderWidget(snap webdm.Package) scopes.PreviewWidget {
	header := scopes.NewPreviewWidget("header", "header")

	header.AddAttributeMapping("title", "title")
	header.AddAttributeMapping("subtitle", "subtitle")
	header.AddAttributeMapping("mascot", "art")

	// According to the store design, the price is only displayed if the app
	// isn't installed.
	if !snap.Installed() && !snap.Installing() {
		priceAttribute := make(map[string]interface{})
		priceAttribute["value"] = "FREE" // All the snaps are currently free
		header.AddAttributeValue("attributes", []interface{}{priceAttribute})
	}

	return header
}

// packageActionWidget is used to create an action widget to install/uninstall/open a snap.
//
// Parameters:
// snap: webdm.Package representing snap.
//
// Returns:
// - Action preview widget for the given snap.
func packageActionWidget(snap webdm.Package) scopes.PreviewWidget {
	actions := scopes.NewPreviewWidget("actions", "actions")

	// The buttons need to provide the options to open and uninstall if the
	// app is installed. Otherwise, just provide the option to install.
	if snap.Installed() || snap.Installing() {
		openAction := make(map[string]interface{})
		openAction["id"] = "open"
		openAction["label"] = "Open"

		uninstallAction := make(map[string]interface{})
		uninstallAction["id"] = "uninstall"
		uninstallAction["label"] = "Uninstall"

		actions.AddAttributeValue("actions", []interface{}{openAction, uninstallAction})
	} else {
		installAction := make(map[string]interface{})
		installAction["id"] = "install"
		installAction["label"] = "Install"

		actions.AddAttributeValue("actions", []interface{}{installAction})
	}

	return actions
}

// scopeError prints an error to stderr as well as returning an actual error.
// This is used because the errors returned from the scope functions don't seem
// to be handled, logged, or otherwise displayed to the user in any way.
func scopeError(format string, a ...interface{}) error {
	log.Printf(format, a...)
	return fmt.Errorf(format, a...)
}
