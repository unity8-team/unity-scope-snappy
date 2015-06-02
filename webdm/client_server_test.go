package webdm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// FakeResponse holds the response from the server
type FakeResponse struct {
	Package string // Package name
	Message string // The message associated with the package
}

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the webdm client being tested.
	client *Client

	// server is a test HTTP server used to provide fake API responses.
	server *httptest.Server

	// storePackages holds all the packages available in the store
	storePackages []Package
)

// fakeSetup simply sets up a test HTTP server along with a webdm.Client that is
// configured to talk to it.
func fakeSetup() {
	// Test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// webdm client configured to use test server
	client, _ = NewClient(server.URL)
}

// setup sets up a test HTTP server along with a webdm.Client that is configured
// to talk to that test server. It then initializes the store and configures
// handlers for the fake API.
func setup(t *testing.T) {
	fakeSetup()

	// Setup fake API
	initializeStore()
	setupHandlers(t)
}

// setup sets up a test HTTP server along with a webdm.Client that is configured
// to talk to that test server. It then registers a handler for the packages
// API that simply returns a server error (500).
func setupBroken() {
	fakeSetup()

	// Handle anything in the packages API, and return a 500.
	mux.HandleFunc(apiPackagesPath,
		func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "Internal Server Error", 500)
		})
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

// initializeStore sets up the fake package store to hold a few packages
func initializeStore() {
	storePackages = []Package{
		Package{
			Id:           "package1",
			Name:         "package1",
			Origin:       "foo",
			Version:      "0.1",
			Vendor:       "bar",
			Description:  "baz",
			IconUrl:      "http://fake",
			Status:       StatusNotInstalled,
			DownloadSize: 123456,
			Type:         "oem",
		},
		Package{
			Id:           "package2",
			Name:         "package2",
			Origin:       "foo",
			Version:      "0.1",
			Vendor:       "bar",
			Description:  "baz",
			IconUrl:      "http://fake",
			Status:       StatusNotInstalled,
			DownloadSize: 123456,
			Type:         "app",
		},
	}
}

// setupHandlers sets up handlers for listing, querying, installing, and
// uninstalling packages.
func setupHandlers(t *testing.T) {
	mux.HandleFunc(apiPackagesPath,
		func(writer http.ResponseWriter, request *http.Request) {
			packageName := request.URL.Path[len(apiPackagesPath):]

			// If no package name was supplied, then list all packages
			if packageName == "" {
				testMethod(t, request, "GET")

				installedOnly := (request.FormValue("installed_only") == "true")

				handlePackageListRequest(t, installedOnly, writer)
				return
			}

			// Try to emulate what WebDM does if the ID is actually another
			// resource path: return a 404
			if strings.Contains(packageName, "/") {
				http.Error(writer, "404 page not found", http.StatusNotFound)
				return
			}

			switch request.Method {
			case "GET":
				handleQueryRequest(t, writer, packageName)
			case "PUT":
				handleInstallRequest(t, writer, packageName)
			case "DELETE":
				handleUninstallRequest(t, writer, packageName)
			default: // Anything else is an error
				t.Error("Unexpected HTTP method: %s", request.Method)
			}
		})
}

// handlePackageListRequest writes the package list into the response as JSON.
//
// Parameters:
// t: Testing handle for any errors that occur.
// installedOnly: Whether the list should only contain installed packages.
// writer: Response writer into which the response will be written.
func handlePackageListRequest(t *testing.T, installedOnly bool, writer http.ResponseWriter) {
	finishOperations()

	// This is pretty inefficient, but fine for these small sets
	packages := make([]Package, 0)
	for _, thisPackage := range storePackages {
		if !installedOnly || thisPackage.Installed() {
			packages = append(packages, thisPackage)
		}
	}

	writer.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(writer)
	encoder.Encode(packages)

	clearMessages()
}

// handleQueryRequest writes the information for a specific package into the
// response as JSON.
//
// Parameters:
// t: Testing handle for any errors that occur.
// writer: Response writer into which the response will be written.
// packageId: ID of the package whose information is being requested.
func handleQueryRequest(t *testing.T, writer http.ResponseWriter, packageId string) {
	packageAvailable := false
	var thisPackage Package
	var index int
	for index, thisPackage = range storePackages {
		if thisPackage.Id == packageId {
			packageAvailable = true
			break
		}
	}

	encoder := json.NewEncoder(writer)

	if !packageAvailable {
		writer.WriteHeader(http.StatusNotFound)
		encoder.Encode(fmt.Sprintf("snappy package not found %s\n", packageId))
		return
	}

	storePackages[index].Status = finishedOperation(thisPackage.Status)

	writer.WriteHeader(http.StatusOK)
	encoder.Encode(storePackages[index])

	// Clear this package's message (if any) and resolve any pending operations
	storePackages[index].Message = ""
}

// handleInstallRequest sets a given package as installed and always returns
// a 202. The body will contain different things depending on the status of
// the install.
//
// Parameters:
// t: Testing handle for any errors that occur.
// writer: Response writer into which the response will be written.
// packageId: ID of the package to install
func handleInstallRequest(t *testing.T, writer http.ResponseWriter, packageId string) {
	// The real WebDM doesn't seem to care whether or not the package is even
	// available. The API always response with either a 202 or a 400, depending
	// on whether an install request has already come in for this package.
	encoder := json.NewEncoder(writer)

	for index, thisPackage := range storePackages {
		if thisPackage.Id == packageId {
			message := "Accepted"
			if operationPending(thisPackage) {
				message = "Operation in progress"
				writer.WriteHeader(http.StatusBadRequest)
			} else {
				writer.WriteHeader(http.StatusAccepted)
			}

			fakeResponse := &FakeResponse{
				Package: packageId,
				Message: message,
			}

			encoder.Encode(fakeResponse)

			storePackages[index].Status = StatusInstalling
			break
		}
	}
}

// handleUninstallRequest sets a given package as uninstalled
func handleUninstallRequest(t *testing.T, writer http.ResponseWriter, packageId string) {
	// The real WebDM doesn't seem to care whether or not the package is even
	// available. The API always response with either a 202 or a 400, depending
	// on whether an install request has already come in for this package.
	encoder := json.NewEncoder(writer)

	for index, thisPackage := range storePackages {
		if thisPackage.Id == packageId {
			message := "Accepted"
			if operationPending(thisPackage) {
				message = "Operation in progress"
				writer.WriteHeader(http.StatusBadRequest)
			} else {
				writer.WriteHeader(http.StatusAccepted)
			}

			fakeResponse := &FakeResponse{
				Package: packageId,
				Message: message,
			}

			encoder.Encode(fakeResponse)

			storePackages[index].Status = StatusUninstalling
			break
		}
	}
}

// testMethod ensures the HTTP method used is the one expected
func testMethod(t *testing.T, request *http.Request, expected string) {
	if expected != request.Method {
		t.Errorf("Request method was %s, expected %s", request.Method, expected)
	}
}

// operationPending determines if an operation is currently pending on a given
// package.
//
// Parameters:
// snap: The package that will be checked for pending operations.
func operationPending(snap Package) bool {
	return snap.Installing() || snap.Uninstalling()
}

// clearMessages iterates through the store packages, clearing any package
// messages.
func clearMessages() {
	for index, _ := range storePackages {
		storePackages[index].Message = ""
	}
}

// finishOperations iterates through the store packaging, resolving any pending
// operations.
func finishOperations() {
	for index, snap := range storePackages {
		storePackages[index].Status = finishedOperation(snap.Status)
	}
}

// finishedOperation converts a given "pending" status into its resolved state.
//
// Parameters:
// status: Status to convert into a resolved state.
//
// Returns:
// - Resolved status. Note that this is just a copy if the `status` parameter
//   was not in a pending state.
func finishedOperation(status Status) Status {
	switch status {
	case StatusInstalling:
		return StatusInstalled
	case StatusUninstalling:
		return StatusNotInstalled
	default:
		return status
	}
}

// runApiRequest runs the desired API request and decodes the JSON response
// into an interface.
//
// Parameters:
// method: HTTP method to use in request
// path: Path to use for request (relative to server.URL)
// value: Interface into which the response will be decoded.
//
// Returns:
// - Error (nil if none)
func runApiRequest(method string, path string, value interface{}) error {
	baseUrl, err := url.Parse(server.URL)
	if err != nil {
		return fmt.Errorf("Error parsing server URL: %s", err)
	}

	relativeUrl, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("Error parsing relative path: %s", err)
	}

	requestUrl := baseUrl.ResolveReference(relativeUrl).String()

	request, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}

	defer response.Body.Close()

	if value != nil {
		err = json.NewDecoder(response.Body).Decode(value)
		if err != nil {
			return fmt.Errorf("Error decoding response: %s", err)
		}
	}

	return nil
}

// Test that the fake server clears pending operations upon a query like the
// real server
func TestFakeServer_pendingOperationsQuery(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	response := new(FakeResponse)

	// Request installation of "package1"
	err := runApiRequest("PUT", apiPackagesPath+"package1", response)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	if response.Package != "package1" {
		t.Errorf(`Response was for package "%s", expected "package1"`, response.Package)
	}

	if response.Message != "Accepted" {
		t.Errorf(`Response message was "%s", expected "Accepted"`, response.Message)
	}

	// Query "package1." Don't care about the response.
	err = runApiRequest("GET", apiPackagesPath+"package1", nil)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	// Request installation of "package1" again
	err = runApiRequest("PUT", apiPackagesPath+"package1", response)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	if response.Package != "package1" {
		t.Errorf(`Response was for package "%s", expected "package1"`, response.Package)
	}

	// Should be "Accepted" again, since we cleared the pending operation
	// with the query.
	if response.Message != "Accepted" {
		t.Errorf(`Response message was "%s", expected "Accepted"`, response.Message)
	}
}

// Test that the fake server clears pending operations upon a package list like
// the real server
func TestFakeServer_pendingOperationsList(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	response := new(FakeResponse)

	// Request installation of "package1"
	err := runApiRequest("PUT", apiPackagesPath+"package1", response)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	if response.Package != "package1" {
		t.Errorf(`Response was for package "%s", expected "package1"`, response.Package)
	}

	if response.Message != "Accepted" {
		t.Errorf(`Response message was "%s", expected "Accepted"`, response.Message)
	}

	// Request package list. Don't care about the response.
	err = runApiRequest("GET", apiPackagesPath, nil)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	// Request installation of "package1" again
	err = runApiRequest("PUT", apiPackagesPath+"package1", response)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	if response.Package != "package1" {
		t.Errorf(`Response was for package "%s", expected "package1"`, response.Package)
	}

	// Should be "Accepted" again, since we cleared the pending operation
	// with the package list request.
	if response.Message != "Accepted" {
		t.Errorf(`Response message was "%s", expected "Accepted"`, response.Message)
	}
}

// Test that the fake server can deal with two requests to install the same
// package.
func TestFakeServer_twoInstallRequests(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	response := new(FakeResponse)

	// Request installation of "package1"
	err := runApiRequest("PUT", apiPackagesPath+"package1", response)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	if response.Package != "package1" {
		t.Errorf(`Response was for package "%s", expected "package1"`, response.Package)
	}

	if response.Message != "Accepted" {
		t.Errorf(`Response message was "%s", expected "Accepted"`, response.Message)
	}

	// Request installation of "package1" again
	err = runApiRequest("PUT", apiPackagesPath+"package1", response)
	if err != nil {
		t.Errorf("Unexpected error running API request: %s", err)
	}

	if response.Package != "package1" {
		t.Errorf(`Response was for package "%s", expected "package1"`, response.Package)
	}

	if response.Message != "Operation in progress" {
		t.Errorf(`Response message was "%s", expected "Operation in progress"`, response.Message)
	}
}
