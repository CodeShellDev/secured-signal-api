package utils

/*
 * General Functions (utils)
 * Might move Functions into seperate files
 */

import (
	"encoding/json"
	"strings"
)

func StringToArray(sliceStr string) []string {
    if sliceStr == "" {
        return nil
    }

    rawItems := strings.Split(sliceStr, ",")
    items := make([]string, 0, len(rawItems))

    for _, item := range rawItems {
        trimmed := strings.TrimSpace(item)
        if trimmed != "" {
            items = append(items, trimmed)
        }
    }

    return items
}

func GetJsonSafe[T any](jsonStr string) (T, error) {
	var result T

	err := json.Unmarshal([]byte(jsonStr), &result)

	return result, err
}

func GetJson[T any](jsonStr string) (T) {
	var result T

	err := json.Unmarshal([]byte(jsonStr), &result)

	if err != nil {
		// JSON is empty
	}

	return result
}