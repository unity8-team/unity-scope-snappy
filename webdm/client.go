package webdm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// User-agent to use when communicating with webdm API
	defaultUserAgent = "unity-scope-snappy"

	// Where webdm is listening
	defaultApiUrl = "http://127.0.0.1:4200"

	// webdm API path to use to obtain a list of packages
	apiListPackagesPath = "/api/v2/packages"
)

// Unmarshall the Status field in the json into a "Installed" boolean
func (s *Status) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("Status: UnmarshalJSON on nil pointer")
	}
	*s = string(data) == "installed"
	return nil
}

// Client is the main struct allowing for communication with the webdm API.
type Client struct {
	// Actual HTTP client for communicating with the webdm API
	client *http.Client

	// Base URL for API requests
	BaseUrl *url.URL

	// User agent used when communicating with the API
	UserAgent string
}

// NewClient creates a new client for communicating with the webdm API
func NewClient() *Client {
	client := new(Client)
	client.client = http.DefaultClient
	client.UserAgent = defaultUserAgent
	client.BaseUrl, _ = url.Parse(defaultApiUrl)

	return client
}

// GetInstalledPackages sends an API request for a list of installed packages.
//
// Returns:
// - Slice of Packags structs
// - Error (nil of none)
func (client Client) GetInstalledPackages() ([]Package, error) {
	packages, err := client.getPackages(true)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error getting installed packages: %s", err)
	}

	return packages, nil
}

// GetStorePackages sends an API request for a list of all packages in the
// store (including installed packages).
//
// Returns:
// - Slice of Packags structs
// - Error (nil of none)
func (client Client) GetStorePackages() ([]Package, error) {
	packages, err := client.getPackages(false)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error getting store packages: %s", err)
	}

	return packages, nil
}

// getPackages sends a request to the API for a package list.
//
// Parameters:
// installedOnly: Whether the list should only contain installed packages.
//
// Returns:
// - Slice of Package structs
// - Error (nil if none)
func (client Client) getPackages(installedOnly bool) ([]Package, error) {
	data := url.Values{}
	if installedOnly {
		data.Set("installed_only", "true")
	}

	request, err := client.newRequest("GET", apiListPackagesPath, data)
	if err != nil {
		return nil, fmt.Errorf("Error creating API request: %s", err)
	}

	var packages []Package
	_, err = client.do(request, &packages)
	if err != nil {
		return nil, fmt.Errorf("Error making API request: %s", err)
	}

	return packages, nil
}

// newRequest creates an API request.
//
// Parameters:
// method: HTTP method (e.g. GET, POST, etc.)
// path: API path relative to BaseUrl
// query: key-values which will be included in the request URL as a query
//
// Returns:
// - Pointer to created HTTP request
// - Error (nil if none)
func (client *Client) newRequest(method string, path string, query url.Values) (*http.Request, error) {
	relativeUrl, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("Eror parsing path %s: %s", path, err)
	}

	requestUrl := client.BaseUrl.ResolveReference(relativeUrl)

	// Create the request. The current webdm API doesn't support bodies
	request, err := http.NewRequest(method, requestUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	request.URL.RawQuery = query.Encode() // Add the desired URL query

	// Add the configured user agent to the request
	request.Header.Add("User-Agent", client.UserAgent)

	return request, nil
}

// do sends an API request. The response is decoded and marshalled into `value`.
//
// Parameters:
// request: HTTP request representing the API call.
// value: Destination of the decoded API response (if the response does not
//        successfully decode into this type, it will cause an error).
//
// Returns:
// - Pointer to resulting HTTP response (even upon error)
// - Error (nil if none)
func (client *Client) do(request *http.Request, value interface{}) (*http.Response, error) {
	response, err := client.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error making API request: %s", err)
	}

	defer response.Body.Close()

	err = checkResponse(response)
	if err != nil {
		// Return the response in case the caller wants to investigate futher
		return response, fmt.Errorf("Error in API response: %s", err)
	}

	// Assuming we were given a value, decode into it
	if value != nil {
		err = json.NewDecoder(response.Body).Decode(value)
		if err != nil {
			// Still return the response in case the caller is interested
			return response, fmt.Errorf("Error decoding response: %s",
				err)
		}
	}

	return response, nil
}

// checkResponse ensures the server response means it's okay to continue.
//
// Parameters:
// response: Response from the server that will be checked.
//
// Returns:
// - Error (nil if none)
func checkResponse(response *http.Response) error {
	code := response.StatusCode
	if code < 200 || code > 299 {
		return fmt.Errorf("Response was %d", code)
	}

	return nil
}
