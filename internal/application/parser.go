package application

import (
	"log"
	"strings"
)

// Parses passed body and returns a map suitable for pushing into storage.
func (a *Application) parse(body string) map[string]string {
	data := make(map[string]string)

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

			// Parse params into "name:value" string.
			valuesString := strings.Split(strings.Split(line, "{")[1], "}")[0]

			var (
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

			for _, param := range params {
				name += "/" + param
			}
		}

		data[name] = value
	}

	log.Printf("Data parsed: %+v\n", data)

	return data
}
