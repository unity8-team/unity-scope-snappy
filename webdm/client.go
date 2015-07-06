package webdm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	// Where webdm is listening
	DefaultApiUrl = "http://webdm.local:4200"

	// User-agent to use when communicating with webdm API
	defaultUserAgent = "unity-scope-snappy"

	// webdm default icon path
	apiDefaultIconPath = "/public/images/default-package-icon.svg"

	// webdm API path to use for package-specific requests (e.g. list, query,
	// install, etc.)
	apiPackagesPath = "/api/v2/packages/"
)

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
//
// Parameters:
// apiUrl: URL where WebDM is listening (host[:port])
//
// Returns:
// - Pointer to new client
// - Error (nil if none)
func NewClient(apiUrl string) (*Client, error) {
	client := new(Client)
	client.client = http.DefaultClient
	client.UserAgent = defaultUserAgent

	if apiUrl == "" {
		apiUrl = DefaultApiUrl
	}

	var err error
	client.BaseUrl, err = url.Parse(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error parsing URL \"%s\": %s", apiUrl, err)
	}

	return client, nil
}

// GetInstalledPackages sends an API request for a list of installed packages.
//
// Parameters:
// query: Search query for list.
//
// Returns:
// - Slice of Packags structs
// - Error (nil of none)
func (client *Client) GetInstalledPackages(query string) ([]Package, error) {
	packages, err := client.getPackages(query, true)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error getting installed packages: %s", err)
	}

	return packages, nil
}

// GetStorePackages sends an API request for a list of all packages in the
// store (including installed packages).
//
// Parameters:
// query: Search query for list.
//
// Returns:
// - Slice of Packags structs
// - Error (nil of none)
func (client *Client) GetStorePackages(query string) ([]Package, error) {
	packages, err := client.getPackages(query, false)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error getting store packages: %s", err)
	}

	return packages, nil
}

// Query sends an API request for information on a given snappy package.
//
// Parameters:
// packageId: ID of the package of interest (NOT the name).
//
// Returns:
// - Pointer to resulting Package struct.
// - Error (nil if none)
func (client *Client) Query(packageId string) (*Package, error) {
	request, err := client.newRequest("GET", apiPackagesPath+packageId, nil)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error creating API request: %s", err)
	}

	snap := new(Package)
	_, err = client.do(request, snap)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error making API request: %s", err)
	}

	snap.IconUrl = client.fixIconUrl(snap.IconUrl)

	return snap, nil
}

// Install sends an API request for a specific snappy package to be installed.
//
// Parameters:
// packageId: ID of the package to install (NOT the name).
//
// Returns:
// - Error (nil if none). Note that installing a package that is already
//   installed is not considered an error.
func (client *Client) Install(packageId string) error {
	request, err := client.newRequest("PUT", apiPackagesPath+packageId, nil)
	if err != nil {
		return fmt.Errorf("webdm: Error creating API request: %s", err)
	}

	// This could possibly return a 400, which just means that the package is
	// essentially already installed but hasn't yet been refreshed. No need to
	// error out if that's the case.
	response, err := client.do(request, nil)
	if err != nil && response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("webdm: Error making API request: %s", err)
	}

	return nil
}

// Uninstall sends an API request for a specific snappy package to be uninstalled.
//
// Parameters:
// packageId: ID of the package to uninstall (NOT the name).
//
// Returns:
// - Error (nil if none). Note that uninstalling a package that is not installed
//   is not considered an error.
func (client *Client) Uninstall(packageId string) error {
	request, err := client.newRequest("DELETE", apiPackagesPath+packageId, nil)
	if err != nil {
		return fmt.Errorf("webdm: Error creating API request: %s", err)
	}

	// This could possibly return a 400, which just means that the package is
	// essentially already uninstalled but hasn't yet been refreshed. No need to
	// error out if that's the case.
	response, err := client.do(request, nil)
	if err != nil && response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("webdm: Error making API request: %s", err)
	}

	return nil
}

// getPackages sends a request to the API for a package list.
//
// Parameters:
// query: Search query for list.
// installedOnly: Whether the list should only contain installed packages.
//
// Returns:
// - Slice of Package structs
// - Error (nil if none)
func (client *Client) getPackages(query string, installedOnly bool) ([]Package, error) {
	data := url.Values{}
	data.Set("q", query)

	if installedOnly {
		data.Set("installed_only", "true")
	}

	request, err := client.newRequest("GET", apiPackagesPath, data)
	if err != nil {
		return nil, fmt.Errorf("Error creating API request: %s", err)
	}

	var packages []Package
	_, err = client.do(request, &packages)
	if err != nil {
		return nil, fmt.Errorf("Error making API request: %s", err)
	}

	for i, thisPackage := range packages {
		packages[i].IconUrl = client.fixIconUrl(thisPackage.IconUrl)
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
		return nil, err
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
			return response, fmt.Errorf("Error decoding response: %s", err)
		}
	}

	return response, nil
}

// fixIconUrl checks an icon URL to ensure it's pointing to a valid icon.
//
// Invalid icon URLs can occur in two cases:
//
// 1) The app is installed, in which case the API provides an icon URL
//    that is relative to webdm's base URL.
// 2) The icon URL is actually invalid in the store database.
//
// If (1), we can turn it into a valid URL using webdm's base URL. If
// (2), we'll need to use a default icon.
//
// Parameters:
// iconUrlString: Icon URL to be (potentially) fixed.
//
// Returns:
// - Fixed icon URL. Note that if the original didn't need to be fixed, the
//   original is returned.
func (client Client) fixIconUrl(iconUrlString string) string {
	iconUrl, err := url.Parse(iconUrlString)
	if err != nil || (!iconUrl.IsAbs() && iconUrl.Path == "") {
		log.Printf("Invalid icon URL: \"%s\", using default", iconUrlString)

		iconUrl, _ = url.Parse(apiDefaultIconPath)
	}

	// Note that if the icon URL is already absolute, ResolveReference() won't
	// change it.
	return client.BaseUrl.ResolveReference(iconUrl).String()
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
