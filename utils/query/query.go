package query

import (
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/safestrings"
)

func ParseRawQuery(raw string) map[string][]string {
	result := make(map[string][]string)
	pairs := strings.SplitSeq(raw, "&")

	for pair := range pairs {
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)

		key := parts[0]
		val := ""

		if len(parts) == 2 {
			val = parts[1]
		}

		result[key] = append(result[key], val)
	}

	return result
}

func ParseTypedQueryValues(values []string) any {
	raw := values[len(values)-1]

	return safestrings.ToType(raw)
}

func ParseTypedQuery(query string, matchPrefix string) (map[string]any) {
	addedData := map[string]any{}

	queryData := ParseRawQuery(query)

	for key, value := range queryData {
		keyWithoutPrefix, match := strings.CutPrefix(key, matchPrefix)

		if match {
			newValue := ParseTypedQueryValues(value)

			addedData[keyWithoutPrefix] = newValue
		}
	}

	return addedData
}