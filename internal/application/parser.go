package application

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"go.dev.pztrn.name/metricator/pkg/schema"
)

// Parses io.Reader passed and returns a map suitable for pushing into storage.
func (a *Application) parse(r io.Reader) (map[string]schema.Metric, error) {
	data := make(map[string]schema.Metric)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines.
		if line == "" {
			continue
		}

		// Prometheus line contains metric name and metric parameters defined
		// in "{}".
		var (
			name   string
			params []string
		)

		a.logger.Debugln("Analyzing line:", line)

		name = a.getMetricName(line)
		a.logger.Debugln("Got metric name:", name)

		metric, found := data[name]
		if !found {
			a.logger.Debugln("Metric wasn't yet created, creating new structure")

			metric = schema.NewMetric(name, "", "", nil)
		}

		a.logger.Debugf("Got metric to use: %+v\n", metric)

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

			// According to docs HELP and TYPE lines should be printed before actual metric. Do not even
			// report bugs regarding that!
			// Docs: https://github.com/Showmax/prometheus-docs/blob/master/content/docs/instrumenting/exposition_formats.md
			continue
		}

		// Parametrized metrics should have own piece of love - we should
		// add parameters to metric's name. This would also require metrics
		// structure copying.
		if strings.Contains(line, "{") {
			newMetric := metric
			newMetric.Name = newMetric.BaseName

			params = a.getParametersForPrometheusMetric(line)
			for _, param := range params {
				newMetric.Name += "/" + param
			}

			newMetric.Params = params

			metric = newMetric
		}

		metric.Value = a.getMetricValue(line)

		a.logger.Debugf("Got metric: %+v\n", metric)

		data[metric.Name] = metric
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("wasn't able to scan input: %w", err)
	}

	a.logger.Debugf("Data parsed: %+v\n", data)

	return data, nil
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
		// "deeply nested"? I think not. So:
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
