package config

import (
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
)

var transformFuncs = map[string]func(string, any) (string, any) {
	"default": lowercaseTransform,
	"lower": lowercaseTransform,
	"upper": uppercaseTransform,
	"keep":  keepTransform,
}

func keepTransform(key string, value any) (string, any) {
	logger.Info(key)
	return key, value
}

func lowercaseTransform(key string, value any) (string, any) {
	logger.Info(key)
	return strings.ToLower(key), value
}

func uppercaseTransform(key string, value any) (string, any) {
	logger.Info(key)
	return strings.ToUpper(key), value
}