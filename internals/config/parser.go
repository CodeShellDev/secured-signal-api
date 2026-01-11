package config

import (
	"fmt"
	"strings"

	"github.com/codeshelldev/gotl/pkg/configutils"
	"github.com/codeshelldev/gotl/pkg/pretty"
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
		box := pretty.NewAutoBox()
		box.MinWidth = 50
		box.PaddingX = 2
		box.PaddingY = 1

		box.Border.Style = pretty.BorderStyle{
			Color: pretty.Basic(pretty.Yellow),
		}

		box.AddBlock(pretty.Block{
			Align: pretty.AlignCenter,
			Style: pretty.Style{},
			Segments: []pretty.Segment{
				pretty.TextBlockSegment{
					Text: "🚨 Deprecation 🚨",
					Style: pretty.Style{
						Bold: true,
						Foreground: pretty.Basic(pretty.Yellow),
					},
				},
				pretty.InlineSegment{},
				pretty.TextBlockSegment{
					Text: "Please refrain from using",
				},
				pretty.InlineSegment{},
				pretty.TextBlockSegment{
					Text: "`" + source + "`",
					Style: pretty.Style{
						Italic: true,
						Bold: true,
						Background: pretty.Basic(pretty.Red),
					},
				},
				pretty.InlineSegment{},
				pretty.InlineSegment{
					Items: []pretty.Inline{
						pretty.Span{
							Text: "as it has been marked as deprecated",
						},
					},
				},
			},
		})

		fmt.Println(box.Render())
	},
}