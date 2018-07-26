package toggl

import (
	"time"
)

const factor = 1000000000

// Duration is a timeformat used by toggl responses
// It is mainly used for converting this to a
// time.Duration
type Duration int64

// Convert converts Time to time.Duration
func (t *Duration) Convert() time.Duration {
	return time.Duration(*t * factor)
}

func (t *Duration) String() string {
	return t.Convert().String()
}

// DurationFromTimeDuration returns a duration in the toggl format from a
// duration in time.Duration format
func DurationFromTimeDuration(t time.Duration) Duration {
	return Duration(t / factor)
}
