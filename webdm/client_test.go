package webdm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the webdm client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a webdm.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// Test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// webdm client configured to use test server
	client = NewClient()
	client.BaseUrl, _ = url.Parse(server.URL)
}

// setupMockListPackagesApi sets up a test HTTP server along with a webdm.Client
// that is configured to talk to that test server. It also registers a handler
// for the "list packages" API URL which returns a valid package list.
func setupMockListPackagesApi() {
	setup()

	// Setup a handler function to respond to requests to
	// `webdmListPackagesPath`. Our database contains two packages: one
	// installed, one not installed.
	mux.HandleFunc(apiListPackagesPath,
		func(writer http.ResponseWriter, request *http.Request) {
			jsonString := `[
		               {
		                  "id":"package1",
		                  "name":"package1",
		                  "origin":"foo",
		                  "version":"0.1",
		                  "vendor":"bar",
		                  "description":"baz",
		                  "icon":"http://fake",
		                  "status":"installed",
		                  "download_size":123456,
		                  "type":"oem"
		               }`

			if request.FormValue("installed_only") != "true" {
				jsonString += `,
			               {
			                  "id":"package2",
			                  "name":"package2",
			                  "origin":"foo",
			                  "version":"0.1",
			                  "vendor":"bar",
			                  "description":"baz",
			                  "icon":"http://fake",
			                  "status":"uninstalled",
			                  "download_size":123456,
			                  "type":"app"
			               }`
			}

			jsonString += "]"

			fmt.Fprint(writer, jsonString)
		})
}

// setupBrokenMockListPackagesApi sets up a test HTTP server along with a
// webdm.Client that is configured to talk to that test server. It also
// registers a handler for the "list packages" API URL which returns a code 500.
func setupBrokenMockListPackagesApi() {
	setup()

	mux.HandleFunc(apiListPackagesPath,
		func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "Internal Server Error", 500)
		})
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

// testMethod ensures the HTTP method used is the one expected
func testMethod(t *testing.T, request *http.Request, expected string) {
	if expected != request.Method {
		t.Errorf("Request method was %s, expected %s", request.Method, expected)
	}
}

// Test typical NewClient() usage.
func TestNewClient(t *testing.T) {
	client := NewClient()

	expectedUrl, _ := url.Parse(defaultApiUrl)

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

// Test typical newRequest() usage.
func TestNewRequest(t *testing.T) {
	client := NewClient()

	data := url.Values{}
	data.Set("foo", "bar")

	path := "foo"

	expectedUrl, _ := url.Parse(defaultApiUrl)
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
	client := NewClient()

	// Test a bad relative URL first-- ":" is obviously invalid
	_, err := client.newRequest("GET", ":", nil)
	if err == nil {
		t.Error("Expected an error to be returned due to invalid relative URL")
	}

	// Now test a bad base URL
	client.BaseUrl.Host = "%20bar"
	_, err = client.newRequest("GET", "foo", nil)
	if err == nil {
		t.Error("Expected an error to be returned due to invalid base URL")
	}
}

// Test typical do() usage.
func TestDo(t *testing.T) {
	// Run test server
	setup()
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
	setup()
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
	setup()
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
	setup()
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
	client.BaseUrl.Host = "%20bar"
	_, err := client.getPackages(false)
	if err == nil {
		t.Error("Expected error to be returned due to invalid URL")
	}
}

// Test getting packages from an API that returns an error
func TestGetPackages_invalidResponse(t *testing.T) {
	// Run test server
	setupBrokenMockListPackagesApi()
	defer teardown()

	_, err := client.getPackages(false)
	if err == nil {
		t.Error("Expected error to be returned due to broken server")
	}
}

// Test querying for only installed packages
func TestGetInstalledPackages(t *testing.T) {
	// Run test server
	setupMockListPackagesApi()
	defer teardown()

	packages, err := client.GetInstalledPackages()
	if err != nil {
		t.Log("Error: ", err)
	}

	// Expect that only one package will be returned, since only one is
	// installed.
	if len(packages) != 1 {
		// Get out now so we don't use a potentially invalid index later
		t.Fatalf("Number of packages: %d, expected: 1", len(packages))
	}

	package1 := Package{
		Id:           "package1",
		Name:         "package1",
		Origin:       "foo",
		Version:      "0.1",
		Vendor:       "bar",
		Description:  "baz",
		IconUrl:      "http://fake",
		Installed:    true,
		DownloadSize: 123456,
		Type:         "oem",
	}

	if !reflect.DeepEqual(packages[0], package1) {
		t.Error("\"package1\" should have been the only package in response")
	}
}

// Test querying for only installed packages with a server that returns an error
func TestGetInstalledPackages_brokenServer(t *testing.T) {
	// Run test server
	setupBrokenMockListPackagesApi()
	defer teardown()

	_, err := client.GetInstalledPackages()
	if err == nil {
		t.Log("Expected error to be returned due to broken server")
	}
}

// Test querying for all packages
func TestGetStorePackages(t *testing.T) {
	// Run test server
	setupMockListPackagesApi()
	defer teardown()

	packages, err := client.GetStorePackages()
	if err != nil {
		t.Log("Error: ", err)
	}

	if len(packages) != 2 {
		t.Errorf("Number of packages: %d, expected: 2", len(packages))
	}

	package1 := Package{
		Id:           "package1",
		Name:         "package1",
		Origin:       "foo",
		Version:      "0.1",
		Vendor:       "bar",
		Description:  "baz",
		IconUrl:      "http://fake",
		Installed:    true,
		DownloadSize: 123456,
		Type:         "oem",
	}

	package2 := Package{
		Id:           "package2",
		Name:         "package2",
		Origin:       "foo",
		Version:      "0.1",
		Vendor:       "bar",
		Description:  "baz",
		IconUrl:      "http://fake",
		Installed:    false,
		DownloadSize: 123456,
		Type:         "app",
	}

	foundPackage1 := false
	foundPackage2 := false

	// Order is not enforced on the result, so we need to iterate through the
	// slice checking each item.
	for _, thisPackage := range packages {
		if reflect.DeepEqual(thisPackage, package1) {
			foundPackage1 = true
		}

		if reflect.DeepEqual(thisPackage, package2) {
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
	setupBrokenMockListPackagesApi()
	defer teardown()

	_, err := client.GetStorePackages()
	if err == nil {
		t.Log("Expected error to be returned due to broken server")
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
		client = NewClient()
		client.BaseUrl, _ = url.Parse(test.baseUrl)
		fixedUrl := client.fixIconUrl(test.iconUrl)

		if fixedUrl != test.expectedIconUrl {
			t.Errorf("Test case %d: Fixed url was %s, expected %s", i,
				fixedUrl,
				test.expectedIconUrl)
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
				t.Errorf("Test case %d: Expected %d to be acceptable", i,
					test.responseCode)
			}
		} else {
			if checkError == nil {
				t.Errorf("Test case %d: Expected %d to cause an error", i,
					test.responseCode)
			}
		}
	}
}

// Test UnmarshalJSON dies when called on nil object
func TestStatusUnmarshalNil(t *testing.T) {
	var nilStatus *Status = nil
	err := nilStatus.UnmarshalJSON([]byte("this is fake json"))
	if err == nil {
		t.Error("Expected error when Status is nil")
	}
}
