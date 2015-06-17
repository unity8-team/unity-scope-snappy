package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"log"
	"strconv"
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
		"art" : {
			"field": "art"
		},
		"subtitle": "subtitle"
        }
}`

// Scope is the struct representing the scope itself.
type Scope struct {
	webdmClient *webdm.Client
}

// NewScope creates a new Scope using a specific WebDM API URL.
//
// Parameters:
// webdmApiUrl: URL where WebDM is listening.
//
// Returns:
// - Pointer to scope (nil if error).
// - Error (nil if none).
func NewScope(webdmApiUrl string) (*Scope, error) {
	scope := new(Scope)
	var err error
	scope.webdmClient, err = webdm.NewClient(webdmApiUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to create WebDM client: %s", err)
	}

	return scope, nil
}

func (scope *Scope) SetScopeBase(base *scopes.ScopeBase) {
	// Do nothing
}

func (scope Scope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	createDepartments(query, reply)

	packages, err := getPackageList(scope.webdmClient, query.DepartmentID())
	if err != nil {
		return scopeError("unity-scope-snappy: Unable to get package list: %s", err)
	}

	var category *scopes.Category
	if query.DepartmentID() == "installed" {
		category = reply.RegisterCategory("installed_packages", "Installed Packages", "", layout)
	} else {
		category = reply.RegisterCategory("store_packages", "Store Packages", "", layout)
	}

	for _, thisPackage := range packages {
		result := packageResult(category, thisPackage)

		if reply.Push(result) != nil {
			// If the push fails, the query was cancelled. No need to continue.
			return nil
		}
	}

	return nil
}

func (scope Scope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
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

	preview, err := NewPreview(*snap)
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
	// Obtain the ID for this action from the string
	intActionId, err := strconv.Atoi(actionId)
	if err != nil {
		return nil, scopeError(`unity-scope-snappy: Invalid action ID: "%s"`, actionId)
	}

	// Get the action runner corresponding to this action
	actionRunner, err := NewActionRunner(ActionId(intActionId))
	if err != nil {
		return nil, scopeError(`unity-scope-snappy: Unable to handle action "%s": %s`, actionId, err)
	}

	// Obtain the ID for the specific package
	var snapId string
	err = result.Get("id", &snapId)
	if err != nil {
		return nil, scopeError(`unity-scope-snappy: Unable to retrieve ID for package "%s": %s`, result.Title(), err)
	}

	response, err := actionRunner.Run(scope.webdmClient, snapId)
	if err != nil {
		err = scopeError(`unity-scope-snappy: Error handling action "%s": %s`, actionId, err)
	}

	return response, err
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
func packageResult(category *scopes.Category, snap webdm.Package) *scopes.CategorisedResult {
	result := scopes.NewCategorisedResult(category)

	result.SetTitle(snap.Name)
	result.Set("subtitle", snap.Vendor)
	result.SetURI("snappy:" + snap.Id)
	result.SetArt(snap.IconUrl)
	result.Set("id", snap.Id)

	return result
}

// getPackageList is used to obtain a package list for a specific department.
//
// Parameters:
// webdmClient: Configured WebDM Client to obtain package list.
// department: The department whose packages should be listed.
//
// Returns:
// - List of WebDM Package structs
// - Error (nil if none)
func getPackageList(packageManager PackageManager, department string) ([]webdm.Package, error) {
	var packages []webdm.Package
	var err error

	switch department {
	case "installed":
		packages, err = packageManager.GetInstalledPackages()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve installed packages: %s", err)
		}

	default:
		packages, err = packageManager.GetStorePackages()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve store packages: %s", err)
		}
	}

	return packages, nil
}

// scopeError prints an error to stderr as well as returning an actual error.
// This is used because the errors returned from the scope functions don't seem
// to be handled, logged, or otherwise displayed to the user in any way.
func scopeError(format string, a ...interface{}) error {
	log.Printf(format, a...)
	return fmt.Errorf(format, a...)
}
