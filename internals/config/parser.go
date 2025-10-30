package config

import (
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/logger"
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
	logger.Dev("LOWER: ", key)
	return strings.ToLower(key), value
}

func uppercaseTransform(key string, value any) (string, any) {
	logger.Dev("UPPER: ", key)
	return strings.ToLower(key), value
}