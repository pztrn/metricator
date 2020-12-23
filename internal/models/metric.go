package models

// Metric is a generic metric structure.
type Metric struct {
	// BaseName is a metric's base name, used for constructing name.
	BaseName string
	// Name is a metric name.
	Name string
	// Description is a metric description from HELP line.
	Description string
	// Type is a metric type from TYPE line.
	Type string
	// Value is a metric value.
	Value string
	// Params is an additional params which are placed inside "{}".
	Params []string
}

// NewMetric creates new structure for storing single metric data.
func NewMetric(name, mType, description string, params []string) Metric {
	m := Metric{
		BaseName:    name,
		Name:        name,
		Description: description,
		Type:        mType,
		Params:      params,
		Value:       "",
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
