package toggl

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Report is a representation of a toggl report
type Report struct {
	TotalGrand Duration     `json:"total_grand"`
	Error      *ReportError `json:"error"`
}

func createReportRequest(apiKey string, userAgent string, workspaceID int64, projectIDs string, sinceYear int) (*http.Request, error) {
	url, err := url.Parse("https://toggl.com/reports/api/v2/summary")
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()
	queryParams.Set("user_agent", userAgent)
	queryParams.Set("workspace_id", strconv.FormatInt(workspaceID, 10))
	queryParams.Set("since", fmt.Sprintf("%d-01-01", sinceYear))
	if sinceYear < time.Now().Year() {
		queryParams.Set("until", fmt.Sprintf("%d-12-31", sinceYear))
	}
	queryParams.Set("project_ids", projectIDs)
	url.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return req, err
	}
	req.SetBasicAuth(apiKey, "api_token")
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// ReportError comes up when requesting a report failed.
type ReportError struct {
	Message string `json:"message"`
	Tip     string `json:"tip"`
	Code    int    `json:"code"`
}
