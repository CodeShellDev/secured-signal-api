package prettylog

import (
	"regexp"
	"runtime/debug"
	"strings"
)

func prettyStack(from, under, n int) string {
	lines := strings.Split(string(debug.Stack()), "\n")

	if len(lines) >= 7 {
		lines = lines[7:]
	}

	frameSize := 2

	totalFrames := len(lines) / frameSize

	startFrame := under
	endFrame := totalFrames - from

	startFrame = max(0, startFrame)
	endFrame = min(endFrame, totalFrames)
	startFrame = min(startFrame, endFrame)

	if n > 0 && startFrame + n < endFrame {
		endFrame = startFrame + n
	}

	lower := startFrame * frameSize
	upper := endFrame * frameSize

	lines = lines[lower:upper]

	hexCodeRe := regexp.MustCompile(`\+?0x[0-9a-fA-F]+`)
	pathRe := regexp.MustCompile(`^(.*)/([^/]+)$`)
	lineRe := regexp.MustCompile(`:([0-9]+)\s*`)
	funcRe := regexp.MustCompile(`\.(\w+)(\([^)]*\))\s*$`)

	for i := 0; i < len(lines); i++ {
		res := strings.TrimSpace(lines[i])

		res = hexCodeRe.ReplaceAllString(res, "")
		res = pathRe.ReplaceAllString(res, "{fg=gray}$1/{/}$2")
		res = lineRe.ReplaceAllString(res, ":{fg=blue}$1{/}")
		res = funcRe.ReplaceAllString(res, ".{fg=bright_yellow}$1{/}{fg=gray}(){/}")

		lines[i] = res
	}

	return strings.Join(lines, "\n")
}

func highlightKeywords(s string) string {
	keywords := map[string]string{
		"error":  	"fg=red",
		"panic":  	"b,fg=red",
		"failed": 	"i,fg=red",
		"map":    	"fg=bright_green",
		"struct": 	"fg=bright_blue",
		"slice":  	"fg=bright_green",
		"int":    	"fg=blue",
		"[0-9]+":   "fg=blue",
		"string": 	"fg=green",
		"nil": 		"bg=gray",
	}

	for k, style := range keywords {
		re := regexp.MustCompile(`\b(` + k + `)\b`)
		s = re.ReplaceAllString(s, "{"+style+"}$1{/}")
	}

	return s
}

func transformQuotes(str string, fn func(open, close, inside string) string) string {
	var builder strings.Builder
	n := len(str)

	closeMap := map[rune]rune{
		'`': '`',
		'\'': '\'',
		'"': '"',
	}

	for i := 0; i < n; {
		char := str[i]

		closeCh, ok := closeMap[rune(char)]
		if !ok {
			builder.WriteByte(char)
			i++
			continue
		}

		open := char
		i++

		start := i
		for i < n && rune(str[i]) != closeCh {
			i++
		}

		inside := str[start:i]

		if i < n && rune(str[i]) == closeCh {
			close := str[i]
			i++

			builder.WriteString(fn(string(open), string(close), inside))
		} else {
			builder.WriteByte(open)
			builder.WriteString(inside)

			break
		}
	}

	return builder.String()
}