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

	// Setup a handler function to respond to requests to
	// `webdmListPackagesPath`. Our database contains two packages: one
	// installed, one not installed.
	mux.HandleFunc(webdmListPackagesPath,
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

	// webdm client configured to use test server
	client = NewClient()
	client.BaseUrl, _ = url.Parse(server.URL)
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

	if client.BaseUrl.Scheme != defaultWebdmScheme {
		t.Errorf("NewClient BaseUrl.Scheme was %s, expected %s",
			client.BaseUrl.Scheme,
			defaultWebdmScheme)
	}

	if client.BaseUrl.Host != defaultWebdmHost {
		t.Errorf("NewClient BaseUrl.Host was %s, expected %s",
			client.BaseUrl.Host,
			defaultWebdmHost)
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

	expectedUrl := url.URL{
		Scheme:   defaultWebdmScheme,
		Host:     defaultWebdmHost,
		Path:     "foo",
		RawQuery: data.Encode(),
	}

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
func TestNewRequest_badURL(t *testing.T) {
	client := NewClient()

	// ":" is obviously an invalid URL
	_, err := client.newRequest("GET", ":", nil)
	if err == nil {
		t.Error("Expected an error to be returned")
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

// Test querying for only installed packages
func TestGetPackages_onlyInstalled(t *testing.T) {
	// Run test server
	setup()
	defer teardown()

	packages, err := client.GetPackages(true)
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

// Test querying for all packages
func TestGetStorePackages(t *testing.T) {
	// Run test server
	setup()
	defer teardown()

	packages, err := client.GetPackages(false)
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

// Test the response checker is only good with 2xx values
func TestCheckResponse(t *testing.T) {
	response := new(http.Response)

	// Test OK
	response.StatusCode = 200
	if checkResponse(response) != nil {
		t.Error("Expected 200 to be acceptable")
	}

	// Test fulfilled request
	response.StatusCode = 226
	if checkResponse(response) != nil {
		t.Error("Expected 226 to be acceptable")
	}

	// Test URI too long
	response.StatusCode = 122
	if checkResponse(response) == nil {
		t.Error("Expected 122 to cause an error")
	}

	// Test a redirection
	response.StatusCode = 300
	if checkResponse(response) == nil {
		t.Error("Expected 300 to cause an error")
	}
}
