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
		deprecationHandler(source, target)
	},
}

var deprecationHandledMap = map[string]bool{}

func deprecationHandler(source string, target configutils.TransformTarget) {
	handled, _ := deprecationHandledMap[source]

	if handled {
		return
	}

	deprecationHandledMap[source] = true

	msgMap := configutils.ParseTag(target.Source.Tag.Get("deprecation"))

	box := pretty.NewAutoBox()
	box.MinWidth = 60
	box.PaddingX = 2
	box.PaddingY = 1

	box.Border.Style = pretty.BorderStyle{
		Color: pretty.Basic(pretty.Yellow),
	}

	messageParts := strings.Split(configutils.GetValueWithSource(source, target.Parent, msgMap), "\n")
	messageSegments := []pretty.Segment{}

	for _, part := range messageParts {
		messageSegments = append(messageSegments, pretty.StyledTextBlockSegment{
			Raw: part,
		})
	}

	atRoot := !strings.Contains(source, ".")
	refrainPrefix := ""
	refrainSuffix := ""

	if atRoot {
		refrainPrefix = "⇧ "
		refrainSuffix = " (at root)"
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
			pretty.InlineSegment{
				Items: []pretty.Inline{
					pretty.Span{
						Text: refrainPrefix,
						Style: pretty.Style{
							Bold: true,
							Foreground: pretty.Basic(pretty.BrightWhite),
						},
					},
					pretty.Span{
						Text: "`" + source + "`",
						Style: pretty.Style{
							Italic: true,
							Bold: true,
							Background: pretty.Basic(pretty.Red),
						},
					},
					pretty.Span{
						Text: refrainSuffix,
					},
				},
			},
			pretty.InlineSegment{},
			pretty.InlineSegment{
				Items: []pretty.Inline{
					pretty.Span{
						Text: "as it has been marked as ",
					},
					pretty.Span{
						Text: "deprecated",
						Style: pretty.Style{
							Bold: true,
						},
					},
					pretty.Span{
						Text: ":",
					},
				},
			},
			pretty.InlineSegment{},
		},
	})

	box.AddBlock(pretty.Block{
		Segments: messageSegments,
		Align: pretty.AlignCenter,
		Style: pretty.Style{
			Background: pretty.Basic(pretty.BrightBlack),
		},
	})

	box.AddBlock(pretty.Block{
		Align: pretty.AlignCenter,
		Segments: []pretty.Segment{
			pretty.InlineSegment{},
			pretty.TextBlockSegment{
				Text: "Update your config before the next update,\nwhere it will be removed for good",
				Style: pretty.Style{
					Italic: true,
				},
			},
		},
	})

	fmt.Println(box.Render())
}