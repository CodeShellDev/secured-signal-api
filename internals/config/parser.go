package config

import (
	"strings"

	"github.com/codeshelldev/gotl/pkg/configutils"
	"github.com/codeshelldev/secured-signal-api/utils/deprecation"
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
	"changed": func(source string, target configutils.TransformTarget) {
		changedHandler(source, target)
	},
}

func deprecationHandler(source string, target configutils.TransformTarget) {
	msgMap := configutils.ParseTag(target.Source.Tag.Get("deprecation"))

	message := configutils.GetValueWithSource(source, target.Parent, msgMap)

	atRoot := !strings.Contains(source, ".")
	usingPrefix := ""
	usingSuffix := ""

	if atRoot {
		usingPrefix = "⇧ "
		usingSuffix = " (at root)"
	}

	deprecation.Warn(source, deprecation.DeprecationMessage{
		Using: "{b,fg=bright_white}" + usingPrefix + "{/}{b,i,bg=red}`" + source + "`{/}" + usingSuffix,
		Message: message,
		Fix: "",
		Note: "\n{i}Update your config as {b}soon{/} as possible{/}",
	})
}

func brokenHandler(source string, target configutils.TransformTarget) {
	msgMap := configutils.ParseTag(target.Source.Tag.Get("breaking"))

	message := configutils.GetValueWithSource(source, target.Parent, msgMap)

	atRoot := !strings.Contains(source, ".")
	usingPrefix := ""
	usingSuffix := ""

	if atRoot {
		usingPrefix = "⇧ "
		usingSuffix = " (at root)"
	}

	deprecation.Error(source, deprecation.DeprecationMessage{
		Using: "{b,fg=bright_white}" + usingPrefix + "{/}{b,i,bg=red}`" + source + "`{/}" + usingSuffix,
		Message: message,
		Fix: "",
		Note: "\n{i}Update your config {b,fg=red}NOW!{/}{/}",
	})
}

func changedHandler(source string, target configutils.TransformTarget) {
	msgMap := configutils.ParseTag(target.Source.Tag.Get("changing"))

	message := configutils.GetValueWithSource(source, target.Parent, msgMap)

	atRoot := !strings.Contains(source, ".")
	usingPrefix := ""
	usingSuffix := ""

	if atRoot {
		usingPrefix = "⇧ "
		usingSuffix = " (at root)"
	}

	deprecation.Info(source, deprecation.DeprecationMessage{
		Using: "{b,fg=bright_white}" + usingPrefix + "{/}{b,i,bg=blue}`" + source + "`{/}" + usingSuffix,
		Message: message,
		Fix: "",
		Note: "\n{i}Please {b,fg=green}verify{/} if your config has been updated{/}",
	})
}