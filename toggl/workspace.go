package toggl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Workspace is a collection of projects and users
// with a paid or free plan of toggl.
type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Text converts a workspace to a string format that's
// suitable for output on a terminal
func (w *Workspace) Text() string {
	return fmt.Sprintf("ID:\t%d\nName:\t%s\n", w.ID, w.Name)
}

// Workspaces reflects the workspaces api endpoint
type Workspaces struct {
	Client *Client // A toggl client
	Path   string  // the path to be added to the toggl api endpoint
}

// Get all workspaces that can be accessed with the credentials
func (w *Workspaces) Get() ([]Workspace, error) {
	req, err := w.createWorkspaceRequest()
	if err != nil {
		return nil, fmt.Errorf("Creating GET request failed: %s", err)
	}
	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Requesting workspaces failed: %s", err)
	}
	defer func() {
		if deferredErr := resp.Body.Close(); deferredErr != nil {
			log.Printf("Closing HTTP Body failed: %s", deferredErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with http status: %s", resp.Status)
	}

	workspaces := []Workspace{}
	buf := &bytes.Buffer{}
	tr := io.TeeReader(resp.Body, buf)
	err = json.NewDecoder(tr).Decode(&workspaces)
	if err != nil {
		return nil, UnmarshalingError{
			Err:      fmt.Errorf("Unmarshaling failed: %s", err),
			Response: resp,
			Data:     buf,
		}
	}
	return workspaces, nil
}

// UnmarshalingError contains the data that failed to unmarshal
type UnmarshalingError struct {
	Err      error
	Response *http.Response
	Data     *bytes.Buffer
}

func (ue UnmarshalingError) Error() string {
	return ue.Err.Error()
}

func (w *Workspaces) createWorkspaceRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", w.Client.Endpoint, w.Path), nil)
	if err != nil {
		return req, err
	}
	req.SetBasicAuth(w.Client.APIKey, "api_token")
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}
