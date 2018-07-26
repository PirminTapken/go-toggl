package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

// GetDurationInput describes the parameters necessary
// to get a duration
type GetDurationInput struct {
	ProjectIDs string
	SinceYear  int
}

// GetDuration uses the configured client to fetch
// the duration of the hours worked this year.
func GetDuration(client *Client, input GetDurationInput) (float64, error) {
	if input.SinceYear < 2006 || input.SinceYear > 2030 {
		return 0, errors.New("SinceYear must be between 2006 and 2030")
	}
	duration := float64(0)
	for y := input.SinceYear; y <= time.Now().Year(); y++ {
		thisDuration, err := GetDurationForYear(client, input.ProjectIDs, y)
		if err != nil {
			return 0, err
		}
		duration += thisDuration
	}
	return duration, nil
}

// GetDurationForYear returns the hours worked in one year
func GetDurationForYear(client *Client, projectIDs string, year int) (float64, error) {
	req, err := client.CreateReportRequest(projectIDs, year)
	if err != nil {
		return 0, fmt.Errorf("creating report request failed: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("getting report failed: %s", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response data failed: %s", err)
	}
	var report Report
	err = json.Unmarshal(data, &report)
	if err != nil {
		return 0, fmt.Errorf("parsing the report failed: %s", err)
	}
	if report.Error != nil {
		return 0, fmt.Errorf("encountered report error %q, tip: %s", report.Error.Message, report.Error.Tip)
	}
	return report.TotalGrand.Convert().Hours(), nil
}
