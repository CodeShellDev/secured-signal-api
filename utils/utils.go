package utils

/*
 * General Functions (utils)
 * Might move Functions into seperate files
 */

import (
	"encoding/json"
	"regexp"
	"strings"
)

func StringToArray(sliceStr string) []string {
	if sliceStr == "" {
		return nil
	}

	re, err := regexp.Compile(`\s+`)

	if err != nil {
		return nil
	}

	normalized := re.ReplaceAllString(sliceStr, "")

	tokens := strings.Split(normalized, ",")

	return tokens
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