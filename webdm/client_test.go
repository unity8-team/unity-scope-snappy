package webdm

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// Test typical NewClient() usage.
func TestNewClient(t *testing.T) {
	client, _ := NewClient("")

	expectedUrl, _ := url.Parse(DefaultApiUrl)

	if client.BaseUrl.Scheme != expectedUrl.Scheme {
		t.Errorf("NewClient BaseUrl.Scheme was %s, expected %s",
			client.BaseUrl.Scheme,
			expectedUrl.Scheme)
	}

	if client.BaseUrl.Host != expectedUrl.Host {
		t.Errorf("NewClient BaseUrl.Host was %s, expected %s",
			client.BaseUrl.Host,
			expectedUrl.Host)
	}

	if client.UserAgent != defaultUserAgent {
		t.Errorf("NewClient UserAgent was %s, expected %s", client.UserAgent,
			defaultUserAgent)
	}
}

// Test NewClient() with an invalid API URL.
func TestNewClient_invalidApiUrl(t *testing.T) {
	_, err := NewClient(":")
	if err == nil {
		t.Error("Expected an error to be returned due to invalid URL")
	}
}

// Test typical newRequest() usage.
func TestNewRequest(t *testing.T) {
	client, _ := NewClient("")

	data := url.Values{}
	data.Set("foo", "bar")

	path := "foo"

	expectedUrl, _ := url.Parse(DefaultApiUrl)
	expectedUrl.Path = path
	expectedUrl.RawQuery = data.Encode()

	request, _ := client.newRequest("GET", path, data)

	// Ensure URL was expanded correctly, including query data
	if request.URL.String() != expectedUrl.String() {
		t.Errorf("NewRequest URL was %s, expected %s", request.URL.String(),
			expectedUrl.String())
	}

	// Ensure form can access query data
	if request.FormValue("foo") != "bar" {
		t.Errorf("NewRequest form should include \"foo\"")
	}

	// Ensure the user agent is correct
	if request.Header.Get("User-Agent") != client.UserAgent {
		t.Errorf("NewRequest User-Agent was %s, expected %s",
			request.Header.Get("User-Agent"),
			client.UserAgent)
	}
}

// Test handling of a bad URL
func TestNewRequest_badUrl(t *testing.T) {
	client, _ := NewClient("")

	// Test a bad relative URL first-- ":" is obviously invalid
	_, err := client.newRequest("GET", ":", nil)
	if err == nil {
		t.Error("Expected an error to be returned due to invalid relative URL")
	}

	// Now test a bad base URL
	client.BaseUrl.Host = "%20"
	_, err = client.newRequest("GET", "foo", nil)
	if err == nil {
		t.Error("Expected an error to be returned due to invalid base URL")
	}
}

// Test typical do() usage.
func TestDo(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	type Body struct {
		Foo string
	}

	// Setup a handler function to respond to requests to the root url
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		fmt.Fprint(writer, `{"Foo":"bar"}`)
	})

	request, _ := client.newRequest("GET", "/", nil)
	body := new(Body)
	client.do(request, body)

	expected := &Body{Foo: "bar"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body was %v, expected %v", body, expected)
	}
}

// Test decoding into incorrect JSON
func TestDo_incorrectJson(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Expect an int, even though the server responds with a string
	type BadBody struct {
		Foo int
	}

	// Setup a handler function to respond to requests to the root url
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		fmt.Fprint(writer, `{"Foo":"bar"}`) // Respond with a string
	})

	request, _ := client.newRequest("GET", "/", nil)
	body := new(BadBody)
	_, err := client.do(request, body)

	if err == nil {
		t.Errorf("Expected an error due to parsing a string into a int")
	}
}

// Test handling of an HTTP 404 error.
func TestDo_httpError(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Setup a handler function to respond to requests to the root url
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Not Found", 404)
	})

	request, _ := client.newRequest("GET", "/", nil)
	_, err := client.do(request, nil)

	if err == nil {
		t.Error("Expected an HTTP 404 error")
	}
}

// Test handling of an infinite redirect loop.
func TestDo_redirectLoop(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Setup a handler function to respond to requests to the root url
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/", http.StatusFound)
	})

	req, _ := client.newRequest("GET", "/", nil)
	_, err := client.do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned due to redirect loop.")
	}
}

// Test getting packages with an invalid URL
func TestGetPackages_invalidUrl(t *testing.T) {
	client.BaseUrl.Host = "%20"
	_, err := client.getPackages(false)
	if err == nil {
		t.Error("Expected error to be returned due to invalid URL")
	}
}

// Test getting packages from an API that returns an error
func TestGetPackages_invalidResponse(t *testing.T) {
	// Run test server
	setupBroken()
	defer teardown()

	_, err := client.getPackages(false)
	if err == nil {
		t.Error("Expected error to be returned due to broken server")
	}
}

// Test querying for only installed packages
func TestGetInstalledPackages(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Set package1 as "installed"
	for index, snap := range storePackages {
		if snap.Id == "package1" {
			storePackages[index].Status = StatusInstalled
		}
	}

	packages, err := client.GetInstalledPackages()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Expect that only one package will be returned, since only one is
	// installed.
	if len(packages) != 1 {
		// Get out now so we don't use a potentially invalid index later
		t.Fatalf("Number of packages: %d, expected: 1", len(packages))
	}

	if packages[0].Id != "package1" {
		t.Error("\"package1\" should have been the only package in response")
	}
}

// Test querying for only installed packages with a server that returns an error
func TestGetInstalledPackages_brokenServer(t *testing.T) {
	// Run test server
	setupBroken()
	defer teardown()

	_, err := client.GetInstalledPackages()
	if err == nil {
		t.Error("Expected error to be returned due to broken server")
	}
}

// Test querying for all packages
func TestGetStorePackages(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	packages, err := client.GetStorePackages()
	if err != nil {
		t.Error("Error: ", err)
	}

	if len(packages) != 2 {
		t.Errorf("Number of packages: %d, expected: 2", len(packages))
	}

	foundPackage1 := false
	foundPackage2 := false

	// Order is not enforced on the result, so we need to iterate through the
	// slice checking each item.
	for _, thisPackage := range packages {
		if thisPackage.Id == "package1" {
			foundPackage1 = true
		}

		if thisPackage.Id == "package2" {
			foundPackage2 = true
		}
	}

	if !foundPackage1 {
		t.Error("\"package1\" not found within response")
	}

	if !foundPackage2 {
		t.Error("\"package2\" not found within response")
	}
}

// Test querying for all packages with a server that returns an error
func TestGetStorePackages_brokenServer(t *testing.T) {
	// Run test server
	setupBroken()
	defer teardown()

	_, err := client.GetStorePackages()
	if err == nil {
		t.Error("Expected error to be returned due to broken server")
	}
}

// Test typical query usage
func TestQuery(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	snap, err := client.Query("package1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if snap.Id != "package1" {
		t.Error("Expected the response to contain \"package1\"")
	}
}

// Testing querying with an invalid URL
func TestQuery_invalidUrl(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	client.BaseUrl.Host = "%20"
	_, err := client.Query("foo")
	if err == nil {
		t.Error("Expected an error due to invalid URL")
	}
}

// Test querying for a non-existing package
func TestQuery_notFound(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	_, err := client.Query("package_not_found")
	if err == nil {
		t.Error("Expected an error due to unavailable package.")
	}
}

// Test typical package installation
func TestInstall(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Request installation of package "package1"
	err := client.Install("package1")
	if err != nil {
		t.Errorf("Unexpected error while installing: %s", err)
	}

	snap, err := client.Query("package1")
	if err != nil {
		// Make this fatal so we don't dereference NULL later
		t.Fatalf("Unexpected error while querying: %s", err)
	}

	if !snap.Installed() {
		t.Error("Expected package to be installed")
	}
}

// Test installing with an invalid URL
func TestInstall_invalidUrl(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	client.BaseUrl.Host = "%20"
	err := client.Install("foo")
	if err == nil {
		t.Error("Expected an error due to invalid URL")
	}
}

// Unfortunately we can't test installing a non-existing package, since not
// even WebDM seems to care about that. So we'll test trying to install a
// a package with an invalid ID instead.
func TestInstall_notFound(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// The forward slash should invalidate this request
	err := client.Install("foo/bar")
	if err == nil {
		t.Error("Expected an error due to invalid package ID")
	}
}

// Test that our API doesn't complain if asked to install a package that is
// already installed.
func TestInstall_redundantInstall(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Request installation of package "package1"
	err := client.Install("package1")
	if err != nil {
		t.Errorf("Unexpected error while installing: %s", err)
	}

	// Request installation of "package1" again. The WebDM API will complain,
	// but ours should handle it.
	err = client.Install("package1")
	if err != nil {
		t.Errorf("Unexpected error while installing again: %s", err)
	}
}

// Test typical package uninstallation
func TestUninstall(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// Set package1 as "installed"
	for index, snap := range storePackages {
		if snap.Id == "package1" {
			storePackages[index].Status = StatusInstalled
		}
	}

	// Verify that package1 is actually installed
	snap, err := client.Query("package1")
	if err != nil {
		t.Errorf("Unexpected error while querying: %s", err)
	}

	if !snap.Installed() {
		t.Errorf(`Expected "package1" to be installed`)
	}

	// Request that the package be uninstalled
	err = client.Uninstall("package1")
	if err != nil {
		t.Errorf("Unexpected error while uninstalling: %s", err)
	}

	snap, err = client.Query("package1")
	if err != nil {
		t.Errorf("Unexpected error while querying: %s", err)
	}

	if !snap.NotInstalled() {
		t.Error(`Expected "package1" to be uninstalled`)
	}
}

// Test uninstalling with an invalid URL
func TestUninstall_invalidUrl(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	client.BaseUrl.Host = "%20"
	err := client.Uninstall("foo")
	if err == nil {
		t.Error("Expected an error due to invalid URL")
	}
}

// Unfortunately we can't test uninstalling a non-existing package, since not
// even WebDM seems to care about that. So we'll test trying to uninstall a
// a package with an invalid ID instead.
func TestUninstall_notFound(t *testing.T) {
	// Run test server
	setup(t)
	defer teardown()

	// The forward slash should invalidate this request
	err := client.Uninstall("foo/bar")
	if err == nil {
		t.Error("Expected an error due to invalid package ID")
	}
}

// Data for TestFixIconUrl
var fixIconUrlTests = []struct {
	baseUrl         string
	iconUrl         string
	expectedIconUrl string
}{
	{"http://example.com", "http://icon.com/icon", "http://icon.com/icon"},
	{"http://example.com", "/icon", "http://example.com/icon"},
	{"http://example.com", "", "http://example.com" + apiDefaultIconPath},
}

// Test that the icon URL fixer works for our cases
func TestFixIconUrl(t *testing.T) {
	for i, test := range fixIconUrlTests {
		client, _ = NewClient(test.baseUrl)
		fixedUrl := client.fixIconUrl(test.iconUrl)

		if fixedUrl != test.expectedIconUrl {
			t.Errorf("Test case %d: Fixed url was %s, expected %s", i,
				fixedUrl, test.expectedIconUrl)
		}
	}
}

// Data for TestCheckResponse
var checkResponseTests = []struct {
	shouldAccept bool
	responseCode int
}{
	{true, 200},  // OK
	{true, 226},  // Fulfilled request
	{false, 122}, // URI too long
	{false, 300}, // Redirection
}

// Test the response checker is only good with 2xx values
func TestCheckResponse(t *testing.T) {
	for i, test := range checkResponseTests {
		response := &http.Response{StatusCode: test.responseCode}
		checkError := checkResponse(response)
		if test.shouldAccept {
			if checkError != nil {
				t.Errorf("Test case %d: Expected %d to be acceptable, got \"%s\"",
					i, test.responseCode, checkError)
			}
		} else {
			if checkError == nil {
				t.Errorf("Test case %d: Expected %d to cause an error", i,
					test.responseCode)
			}
		}
	}
}
