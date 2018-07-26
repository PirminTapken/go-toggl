package toggl

// Project is a collection of time entries.
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
