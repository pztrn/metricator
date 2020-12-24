package schema

// AppsList represents applications list structure from Metricator's API.
type AppsList []string

// IsEmpty returns true if returned applications list is empty.
func (a AppsList) IsEmpty() bool {
	return len(a) == 0
}
