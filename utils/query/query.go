package query

import (
	"regexp"
	"strconv"
	"strings"
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

func tryParseInt(str string) (int, bool) {
	isInt, err := regexp.MatchString(`^\d+$`, str)

	if isInt && err == nil {
		intValue, err := strconv.Atoi(str)

		if err == nil {
			return intValue, true
		}
	}

	return 0, false
}

func ParseTypedQueryValues(values []string) interface{} {
	var result interface{}

	raw := values[0]

	intValue, isInt := tryParseInt(raw)

	if strings.Contains(raw, ",") || (strings.Contains(raw, "[") && strings.Contains(raw, "]")) {
		if strings.Contains(raw, "[") && strings.Contains(raw, "]") {
			escapedStr := strings.ReplaceAll(raw, "[", "")
			escapedStr = strings.ReplaceAll(escapedStr, "]", "")
			raw = escapedStr
		}

		parts := strings.Split(raw, ",")

		var list []interface{}

		for _, part := range parts {
			_intValue, _isInt := tryParseInt(part)

			if _isInt {
				list = append(list, _intValue)
			} else {
				list = append(list, part)
			}
		}
		result = list
	} else if isInt {
		result = intValue
	} else {
		result = raw
	}

	return result
}

func ParseTypedQuery(query string, matchPrefix string) (map[string]interface{}) {
	addedData := map[string]interface{}{}

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