package config

import (
	"strings"
)

var transformFuncs = map[string]func(string, any) (string, any) {
	"default": defaultTransform,
	"lower": lowercaseTransform,
	"upper": uppercaseTransform,
}

func defaultTransform(key string, value any) (string, any) {
	return key, value
}

func lowercaseTransform(key string, value any) (string, any) {
	return strings.ToLower(key), value
}

func uppercaseTransform(key string, value any) (string, any) {
	return strings.ToUpper(key), value
}