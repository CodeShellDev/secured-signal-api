package utils

/*
 * General Functions (utils)
 * Might move Functions into seperate files
 */

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

func StringToArray(sliceStr string) ([]string, error) {
	if sliceStr == "" {
		return []string{}, errors.New("sliceStr is empty")
	}

	re, err := regexp.Compile(`\s+`)

	if err != nil {
		return []string{}, err
	}

	normalized := re.ReplaceAllString(sliceStr, "")

	tokens := strings.Split(normalized, ",")

	return tokens, nil
}

func Contains[T comparable](list []T, item T) (bool){
    for _, match := range list {
        if match == item {
            return true
        }
    }
	
    return false
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