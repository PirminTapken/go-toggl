package toggl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspacesGetUsesCorrectURL(t *testing.T) {
	var expectedURL string
	var expectedResourceString = "some-irritatingly-wrong-endpoint"
	withServer(t,
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			assert.Equal(t, expectedURL,
				fmt.Sprintf("http://%s%s", r.Host, r.RequestURI))
		}),
		func(s *httptest.Server) {
			expectedURL = fmt.Sprintf("%s/%s", s.URL, expectedResourceString)
			t.Logf("Expected URL: %s", expectedURL)

			w := &Workspaces{
				Client: &Client{
					HTTPClient: s.Client(),
					Endpoint:   s.URL,
				},
				Path: expectedResourceString,
			}
			_, err := w.Get()
			// We get an unmarshaling error, but more importantly,
			// we didn't fail in the handler func above!
			// This error is expected as the handler func above doesn't
			// do json as we don't actually need it.
			assert.IsType(t, UnmarshalingError{}, err)
		},
	)
}

func TestWorkspacesGetList(t *testing.T) {
	expectedWorkspaces := []Workspace{{
		ID:   1,
		Name: "First workspace",
	}}

	withServer(t, serveWorkspaces(t, expectedWorkspaces), func(s *httptest.Server) {
		actualWorkspaces, err := (&Client{
			HTTPClient: s.Client(),
			Endpoint:   s.URL,
		}).Workspaces().Get()
		assert.Nil(t, err)
		assert.Equal(t, expectedWorkspaces, actualWorkspaces)
	})
}

func withServer(t *testing.T, h http.Handler, f func(*httptest.Server)) {
	t.Helper()
	called := false
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			called = true
			h.ServeHTTP(rw, r)
		}),
	)
	defer s.Close()
	f(s)
	assert.True(t, called, "Testserver didn't receive any requests")
}

func serveWorkspaces(t *testing.T, ws []Workspace) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(rw).Encode(&ws)
		assert.Nil(t, err)
	})
}
