package toggl

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeEntries(t *testing.T) {
	client := testClient
	startDate := time.Date(2017, 12, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2017, 12, 20, 0, 0, 0, 0, time.UTC)
	expectedStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 8, 25, 13, 0, time.UTC)
	expectedEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 11, 18, 55, 0, time.UTC)
	expectedEntries := []TimeEntry{{
		ID:       1,
		Pid:      1,
		Wid:      1,
		Billable: false,
		Start:    expectedStart,
		Stop:     expectedEnd,
		Duration: DurationFromTimeDuration(expectedEnd.Sub(expectedStart)),
	}}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(rw).Encode(&expectedEntries)
		assert.Nil(t, err)
	}))

	client.Endpoint = server.URL
	client.HTTPClient = server.Client()

	actualEntries, err := client.GetTimeEntries(startDate, endDate)
	assert.Nil(t, err)
	assert.Equal(t, expectedEntries, actualEntries)
}
