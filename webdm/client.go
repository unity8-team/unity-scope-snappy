package webdm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultUserAgent      = "unity-scope-snappy"
	defaultWebdmScheme    = "http"
	defaultWebdmHost      = "127.0.0.1:4200" // where webdm is listening
	webdmListPackagesPath = "/api/v2/packages"
)

// Package contains information about a given package available from the store
// or already installed.
type apiPackage struct {
	Id           string
	Name         string
	Origin       string
	Version      string
	Vendor       string
	Description  string
	Icon         string
	Status       string
	DownloadSize int `json:"download_size"`
	Type         string
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
	client.BaseUrl = &url.URL{
		Scheme: defaultWebdmScheme,
		Host:   defaultWebdmHost,
	}

	return client
}

// GetPackages sends a request to the API for a package list.
//
// Parameters:
// installedOnly: Whether the list should only contain installed packages.
//
// Returns:
// - Slice of Package structs
// - Error (nil if none)
func (client Client) GetPackages(installedOnly bool) ([]Package, error) {
	data := url.Values{}
	if installedOnly {
		data.Set("installed_only", "true")
	}

	request, err := client.newRequest("GET", webdmListPackagesPath, data)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error creating API request: %s", err)
	}

	var apiPackages []apiPackage
	_, err = client.do(request, &apiPackages)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error making API request: %s", err)
	}

	return convertApiResponse(apiPackages), nil
}

// newRequest creates an API request.
//
// Parameters:
// method: HTTP method (e.g. GET, POST, etc.)
// path: API path relative to BaseUrl
// data: key-values which will be included in the request URL
//
// Returns:
// - Pointer to created HTTP request
// - Error (nil if none)
func (client *Client) newRequest(method string, path string,
	data url.Values) (*http.Request, error) {
	relativeUrl, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("webdm: Eror parsing path %s: %s", path, err)
	}

	requestUrl := client.BaseUrl.ResolveReference(relativeUrl)

	// Create the request. The current webdm API doesn't support bodies
	request, err := http.NewRequest(method, requestUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("webdm: Error creating request: %s", err)
	}

	request.URL.RawQuery = data.Encode() // Add the desired URL data

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

// convertApiResponse is a simple helper to convert from the JSON API-specific
// response into our prettied-up library response.
//
// Parameters:
// apiPackages: Slice of apiPackage structs to convert.
//
// Returns:
// Slice of Package structs to give to client.
func convertApiResponse(apiPackages []apiPackage) []Package {
	packages := make([]Package, len(apiPackages))

	for index, thisApiPackage := range apiPackages {
		newPackage := Package{
			Id:           thisApiPackage.Id,
			Name:         thisApiPackage.Name,
			Origin:       thisApiPackage.Origin,
			Version:      thisApiPackage.Version,
			Vendor:       thisApiPackage.Vendor,
			Description:  thisApiPackage.Description,
			IconUrl:      thisApiPackage.Icon,
			Installed:    false,
			DownloadSize: thisApiPackage.DownloadSize,
			Type:         thisApiPackage.Type,
		}

		if thisApiPackage.Status == "installed" {
			newPackage.Installed = true
		}

		packages[index] = newPackage
	}

	return packages
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
