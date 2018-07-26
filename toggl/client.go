package toggl

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Client is a client to the toggl api
type Client struct {
	APIKey     string
	UserAgent  string
	Workspace  int64
	HTTPClient *http.Client
	Endpoint   string
}

// NewClient returns a toggl client with defaults, ready to use.
func NewClient(apiKey, userAgent string) *Client {
	return &Client{
		APIKey:     apiKey,
		UserAgent:  userAgent,
		HTTPClient: http.DefaultClient,
		Endpoint:   DefaultAPIEndpoint,
	}
}

// Workspaces returns the Workspaces endpoint of the API
func (c *Client) Workspaces() *Workspaces {
	return &Workspaces{
		Client: c,
		Path:   DefaultWorkspacesPath,
	}
}

// CreateProjectsRequest creates a request for projects for a
// specific workspace
func (c *Client) CreateProjectsRequest() (*http.Request, error) {
	urlString := "%s/workspaces/%d/projects"
	url, err := url.Parse(fmt.Sprintf(urlString, c.Endpoint, c.Workspace))
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()
	queryParams.Set("user_agent", c.UserAgent)
	url.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return req, err
	}
	req.SetBasicAuth(c.APIKey, "api_token")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// CreateReportRequest creates a request fot toggl taking care of
// configuration from the environment.
func (c *Client) CreateReportRequest(projectIDs string, sinceYear int) (*http.Request, error) {
	if c.APIKey == "" {
		return nil, errors.New("APIKey must be set")
	}
	if c.UserAgent == "" {
		return nil, errors.New("User agent must be set")
	}
	if c.Workspace == 0 {
		return nil, errors.New("Workspace is set to 0")
	}
	return createReportRequest(c.APIKey, c.UserAgent, c.Workspace, projectIDs, sinceYear)
}

// Do wraps http.Client.Do
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// GetTimeEntries fetches time entries in a specific range
func (c *Client) GetTimeEntries(start, end time.Time) ([]TimeEntry, error) {
	req, err := c.createGetTimeEntriesRequest(start, end)
	if err != nil {
		return nil, fmt.Errorf("Creating GetTimeEntries request failed: %s", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Requesting GetTimeEntries failed: %s", err)
	}
	defer func() {
		// We want only to log the error
		if deferredErr := resp.Body.Close(); deferredErr != nil {
			log.Printf("Closing the HTTP body failed: %s", deferredErr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Request failed: %s", resp.Status)
	}

	entries, err := parseGetTimeEntriesResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse GetTimeEntries response: %s", err)
	}

	return entries, nil
}

func (c *Client) createGetTimeEntriesRequest(start, end time.Time) (*http.Request, error) {

	requestURL := c.createGetTimeEntriesRequestURL(start, end)

	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request failed: %s", err)
	}
	req.SetBasicAuth(c.APIKey, "api_token")
	return req, nil
}

func (c *Client) createGetTimeEntriesRequestURL(start, end time.Time) *url.URL {
	requestURL, err := url.Parse(c.Endpoint + "/time_entries")
	if err != nil {
		// This should not happen except somebody messes up the endpoint
		panic(err)
	}

	queryParams := requestURL.Query()
	queryParams.Set("start_date", start.Format(time.RFC3339))
	queryParams.Set("end_date", end.Format(time.RFC3339))

	requestURL.RawQuery = queryParams.Encode()
	return requestURL
}
