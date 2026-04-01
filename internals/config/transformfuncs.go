package config

import (
	"strings"

	"github.com/codeshelldev/gotl/pkg/configutils"
	"github.com/codeshelldev/secured-signal-api/utils/prettylog"
)

var transformFuncs = map[string]func(string, any) (string, any) {
	"default": lowercaseTransform,
	"lower": lowercaseTransform,
	"upper": uppercaseTransform,
	"keep":  keepTransform,
}

func keepTransform(key string, value any) (string, any) {
	return key, value
}

func lowercaseTransform(key string, value any) (string, any) {
	return strings.ToLower(key), value
}

func uppercaseTransform(key string, value any) (string, any) {
	return strings.ToUpper(key), value
}

var onUseFuncs = map[string]func(source string, target configutils.TransformTarget) {
	"deprecated": func(source string, target configutils.TransformTarget) {
		deprecationHandler(source, target)
	},
	"broken": func(source string, target configutils.TransformTarget) {
		brokenHandler(source, target)
	},
}

func deprecationHandler(source string, target configutils.TransformTarget) {
	msgMap := configutils.ParseTag(target.Source.Tag.Get("deprecation"))

	message := configutils.GetValueWithSource(source, target.Parent, msgMap)

	prettylog.Deprecated(source, prettylog.DeprecationMessage{
		Message: message,
	})
}

func brokenHandler(source string, target configutils.TransformTarget) {
	msgMap := configutils.ParseTag(target.Source.Tag.Get("breaking"))

	message := configutils.GetValueWithSource(source, target.Parent, msgMap)

	prettylog.Breaking(source, prettylog.DeprecationMessage{
		Message: message,
	})
}