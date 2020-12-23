package models

// Metric is a generic metric structure.
type Metric struct {
	// Metric name.
	name string
	// HELP data, if present.
	description string
	// Additional parameters, data inside "{}".
	params []string
	// Metric value.
	value string
}

// NewMetric creates new structure for storing single metric data.
func NewMetric(name, description string, params []string) Metric {
	m := Metric{
		name:        name,
		description: description,
		params:      params,
	}

	return m
}

// GetValue returns metric's value.
func (m *Metric) GetValue() string {
	return m.value
}

// SetValue sets value for metric.
func (m *Metric) SetValue(value string) {
	m.value = value
}
