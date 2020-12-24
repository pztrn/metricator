package schema

// Metrics is a metrics collection response.
type Metrics []*Metric

// IsEmpty returns true if returned applications list is empty.
func (m Metrics) IsEmpty() bool {
	return len(m) == 0
}
