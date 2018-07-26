package toggl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// TimeEntry represents a time entry
type TimeEntry struct {
	ID          int       `json:"id"`
	Pid         int       `json:"pid"`
	Wid         int       `json:"wid"`
	Billable    bool      `json:"billable"`
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Duration    Duration  `json:"duration"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

func parseGetTimeEntriesResponse(r *http.Response) ([]TimeEntry, error) {
	buf := &bytes.Buffer{}
	body := io.TeeReader(r.Body, buf)

	entries := make([]TimeEntry, 0)
	err := json.NewDecoder(body).Decode(&entries)
	if err != nil {
		data, ioErr := ioutil.ReadAll(buf)
		if ioErr != nil {
			panic(ioErr)
		}
		log.Printf("HTTP Body: %#v", string(data))
		log.Printf("HTTP Status: %s", r.Status)
		return nil, fmt.Errorf("Decoding time entries failed: %s", err)
	}
	return entries, nil
}
