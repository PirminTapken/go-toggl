package toggl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAPIKey    = "test-api-key"
	testUserAgent = "test-user-agent"
	testWorkspace = int64(123456)
)

var (
	testHTTPClient = http.DefaultClient
	testClient     = &Client{
		APIKey:     testAPIKey,
		UserAgent:  testUserAgent,
		Workspace:  testWorkspace,
		HTTPClient: testHTTPClient,
	}
)

func TestNewClient(t *testing.T) {
	c := NewClient(testAPIKey, testUserAgent)
	assert.NotNil(t, c)
	assert.Equal(t, c.APIKey, testAPIKey)
	assert.Equal(t, c.UserAgent, testUserAgent)
	assert.Equal(t, c.Workspace, int64(0))
	assert.Equal(t, c.HTTPClient, http.DefaultClient)
	assert.Equal(t, c.Endpoint, DefaultAPIEndpoint)
}

func TestCreateProjectsRequest(t *testing.T) {
	req, err := testClient.CreateProjectsRequest()

	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, testUserAgent, req.UserAgent())
	ua := req.URL.Query().Get("user_agent")
	assert.Equal(t, testUserAgent, ua)
}
