package application

import (
	"log"
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
			name, value string
			params      []string
		)

		// Skip empty lines.
		if line == "" {
			continue
		}

		// Check that line isn't commented. We should skip comments for now.
		if strings.HasPrefix(line, "#") {
			continue
		}

		log.Println("Analyzing line:", line)

		// Check if we have parametrized metric. If no - push it to data map.
		if !strings.Contains(line, "{") {
			name = strings.Split(line, " ")[0]
			value = strings.Split(line, " ")[1]
		} else {
			value = strings.Split(line, " ")[1]
			name = strings.Split(line, "{")[0]
			params = a.getParametersForPrometheusMetric(line)

			for _, param := range params {
				name += "/" + param
			}
		}

		metric := models.NewMetric(name, "", params)
		metric.SetValue(value)

		data[name] = metric
	}

	log.Printf("Data parsed: %+v\n", data)

	return data
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
