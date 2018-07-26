package toggl

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReportRequestThisYear(t *testing.T) {
	// Initial variables for this test
	var (
		apiKey      = "someApiKey"
		userAgent   = "someUserAgent"
		workspaceID = int64(1234)
		projectIDs  = "12,34,56"
		username    = apiKey
		password    = "api_token"
		sinceYear   = 2018
	)

	// Expectations
	var (
		expectedHeader = http.Header{
			"Authorization": []string{createAuthorizationHeader(username, password)},
			"Content-Type":  []string{"application/json"},
		}
	)

	expectedURL := mustParseURL(fmt.Sprintf("https://toggl.com/reports/api/v2/summary?project_ids=%s&since=%d-01-01&user_agent=%s&workspace_id=%d",
		url.PathEscape(projectIDs), sinceYear, userAgent, workspaceID))

	request, err := createReportRequest(apiKey, userAgent, workspaceID, projectIDs, sinceYear)
	assert.Nil(t, err)

	assert.Equal(t, expectedHeader, request.Header)
	assert.Equal(t, expectedURL, request.URL)
}

func TestCreateReportRequestLastYears(t *testing.T) {
	// Initial variables for this test
	var (
		apiKey      = "someApiKey"
		userAgent   = "someUserAgent"
		workspaceID = int64(1234)
		projectIDs  = "12,34,56"
		username    = apiKey
		password    = "api_token"
		sinceYear   = 2017
	)

	// Expectations
	var (
		expectedHeader = http.Header{
			"Authorization": []string{createAuthorizationHeader(username, password)},
			"Content-Type":  []string{"application/json"},
		}
	)

	expectedURL := mustParseURL(fmt.Sprintf("https://toggl.com/reports/api/v2/summary?project_ids=%s&since=%d-01-01&until=%d-12-31&user_agent=%s&workspace_id=%d",
		url.PathEscape(projectIDs), sinceYear, sinceYear, userAgent, workspaceID))

	request, err := createReportRequest(apiKey, userAgent, workspaceID, projectIDs, sinceYear)
	assert.Nil(t, err)

	assert.Equal(t, expectedHeader, request.Header)
	assert.Equal(t, expectedURL, request.URL)
}

func createAuthorizationHeader(username, password string) string {
	authString := strings.Join([]string{username, password}, ":")
	encodedAuthentication := base64.StdEncoding.EncodeToString([]byte(authString))
	headerString := strings.Join([]string{"Basic", encodedAuthentication}, " ")
	return headerString
}

func mustParseURL(rawurl string) *url.URL {
	expectedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return expectedURL
}
