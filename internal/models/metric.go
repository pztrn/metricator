package models

// Metric is a generic metric structure.
type Metric struct {
	// Metric name.
	Name string
	// HELP data, if present.
	Description string
	// Additional parameters, data inside "{}".
	Params []string
	// Type is a metric type.
	Type string
	// Metric value.
	Value string
}

// NewMetric creates new structure for storing single metric data.
func NewMetric(name, mType, description string, params []string) Metric {
	m := Metric{
		Name:        name,
		Description: description,
		Type:        mType,
		Params:      params,
	}

	return m
}

// GetValue returns metric's value.
func (m *Metric) GetValue() string {
	return m.Value
}

// SetValue sets value for metric.
func (m *Metric) SetValue(value string) {
	m.Value = value
}
