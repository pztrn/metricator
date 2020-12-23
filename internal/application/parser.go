package application

import (
	"strings"

	"go.dev.pztrn.name/metricator/internal/models"
)

// Parses passed body and returns a map suitable for pushing into storage.
func (a *Application) parse(body string) map[string]models.Metric {
	data := make(map[string]models.Metric)

	// ToDo: switch to bytes buffer and maybe do not read body in caller?
	splittedBody := strings.Split(body, "\n")

	for _, line := range splittedBody {
		// Prometheus line contains metric name and metric parameters defined
		// in "{}".
		var (
			name   string
			params []string
		)

		// Skip empty lines.
		if line == "" {
			continue
		}

		// log.Println("Analyzing line:", line)

		name = a.getMetricName(line)
		metric, found := data[name]
		if !found {
			metric = models.NewMetric(name, "", "", nil)
		}

		// If line is commented - then we have something about metric's description
		// or type. It should be handled in special way - these metric will became
		// "pseudometric" which will be used as template for next iterations. For
		// example if HELP line was parsed first, then TYPE line will be parsed and
		// data will be added to already existing metric. If next line which should
		// represent metric itself contains parameters (e.g. "{instance='host1'}")
		// then next code block will COPY populated with HELP and TYPE pesudometric
		// and put it's value as new metric with different name (metric/instance:host1
		// in this example).
		if strings.HasPrefix(line, "#") {
			switch strings.Split(line, " ")[1] {
			case "HELP":
				metric.Description = a.getMetricDescription(line)
			case "TYPE":
				metric.Type = a.getMetricType(line)
			}

			data[name] = metric

			// According to https://github.com/Showmax/prometheus-docs/blob/master/content/docs/instrumenting/exposition_formats.md
			// HELP and TYPE lines should be printed before actual metric. Do not even
			// report bugs regarding that!
			continue
		}

		// Parametrized metrics should have own piece of love - we should
		// add parameters to metric's name. This would also require metrics
		// structure copying.
		if strings.Contains(line, "{") {
			newMetric := metric

			params = a.getParametersForPrometheusMetric(line)
			for _, param := range params {
				newMetric.Name += "/" + param
			}

			metric = newMetric
			data[metric.Name] = metric
		}

		metric.Value = a.getMetricValue(line)

		// log.Printf("Got metric: %+v\n", metric)

		data[name] = metric
	}

	// log.Printf("Data parsed: %+v\n", data)

	return data
}

// Gets metric description from passed line.
func (a *Application) getMetricDescription(line string) string {
	return strings.Join(strings.Split(line, " ")[3:], " ")
}

// Gets metric name from passed line.
func (a *Application) getMetricName(line string) string {
	var metricNameData string

	if strings.HasPrefix(line, "#") {
		metricNameData = strings.Split(line, " ")[2]
	} else {
		metricNameData = strings.Split(line, " ")[0]
	}

	return strings.Split(metricNameData, "{")[0]
}

// Gets metric type from passed line.
func (a *Application) getMetricType(line string) string {
	return strings.Split(line, " ")[3]
}

// Gets metric value from passed line.
func (a *Application) getMetricValue(line string) string {
	if strings.Contains(line, "}") {
		return strings.Split(line, "} ")[1]
	}

	return strings.Split(line, " ")[1]
}

// Parses passed line and returns a slice of strings with parameters parsed.
func (a *Application) getParametersForPrometheusMetric(line string) []string {
	valuesString := strings.Split(strings.Split(line, "{")[1], "}")[0]

	var (
		params                                                   []string
		paramName, paramValue                                    string
		paramNameFinished, paramValueStarted, paramValueFinished bool
	)

	for _, r := range valuesString {
		if paramValueFinished && string(r) == "," {
			params = append(params, paramName+":"+paramValue)
			paramName, paramValue = "", ""
			paramNameFinished, paramValueStarted, paramValueFinished = false, false, false

			continue
		}

		// Sometimes nestif causes questions, like here. Is code below is
		// "deply nested"? I think not. So:
		// nolint:nestif
		if !paramNameFinished {
			if string(r) != "=" {
				paramName += string(r)

				continue
			} else {
				paramNameFinished = true

				continue
			}
		} else {
			if string(r) == "\"" && !paramValueStarted {
				paramValueStarted = true

				continue
			}

			if paramValueStarted && string(r) != "\"" {
				paramValue += string(r)

				continue
			}

			if paramValueStarted && string(r) == "\"" {
				paramValueFinished = true

				continue
			}
		}
	}

	if paramName != "" && paramValue != "" {
		params = append(params, paramName+":"+paramValue)
	}

	return params
}
