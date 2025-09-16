package jsonutils

import (
	"encoding/json"
	"regexp"
	"strconv"
)

func GetByPath(path string, data any) (any, bool) {
    // Split into parts by `.` and `[]`
    re := regexp.MustCompile(`\.|\[|\]`)

    parts := re.Split(path, -1)

    cleaned := []string{}

    for _, part := range parts {
        if part != "" {
            cleaned = append(cleaned, part)
        }
    }

    current := data

    for _, key := range cleaned {
        switch currentDataType := current.(type) {
            // Case: Dictionary
            case map[string]any:
                value, ok := currentDataType[key]
                if !ok {
                    return nil, false
                }
                current = value

            // Case: Array
            case []any:
                index, err := strconv.Atoi(key)

                if err != nil || index < 0 || index >= len(currentDataType) {
                    return nil, false
                }
                current = currentDataType[index]

            default:
                return nil, false
        }
    }

    return current, true
}

func GetJsonSafe[T any](jsonStr string) (T, error) {
	var result T

	err := json.Unmarshal([]byte(jsonStr), &result)

	return result, err
}

func GetJson[T any](jsonStr string) (T) {
	var result T

	json.Unmarshal([]byte(jsonStr), &result)

	return result
}

func ToJsonSafe[T any](obj T) (string, error) {
	bytes, err := json.Marshal(obj)

	return string(bytes), err
}

func ToJson[T any](obj T) string {
	bytes, _ := json.Marshal(obj)

	return string(bytes)
}