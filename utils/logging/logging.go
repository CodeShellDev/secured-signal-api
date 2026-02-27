package logging

import (
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config"
)

func DefaultTransforms() []func(string)string {
	transforms := []func(string)string{}

	transforms = append(transforms, BeginWithCapital)

	if config.ENV.REDACT_TOKENS {
		transforms = append(transforms, RedactTokens())
	}

	return transforms
}

func Init(level string) {
	options := logger.DefaultOptions()

	logger.InitWith(level, options)
	logger.InitStdLoggerWith(level, options)
}

func Setup() {
	transform := Apply(DefaultTransforms()...)

	logger.Get().SetTransform(transform)
	logger.GetStdLogger().SetTransform(transform)
}

func RedactTokens() func(string) string {
	return RedactWords('*', config.ENV.TOKENS...)
}

func Apply(transforms ...func(content string) string) func(string) string {
	return func(content string) string {
		for _, fn := range transforms {
			content = fn(content)
		}

		return content
	}
}

func BeginWithCapital(content string) string {
	return strings.ToUpper(content[:1]) + content[1:]
}

func Redact(redact string) string {
	if len(redact) <= 4 {
		return strings.Repeat("*", len(redact))
	}

	return string(redact[0]) + strings.Repeat("*", len(redact) - 2) + string(redact[len(redact) - 1])
}

func RedactWords(replaceBy rune, words ...string) func(string) string {
	return func(content string) string {
		for _, word := range words {
			content = strings.ReplaceAll(content, word, "[" + Redact(word) + "]")
		}

		return content
	}
}