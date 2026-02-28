package deprecation

import (
	"fmt"
	"os"

	"github.com/codeshelldev/gotl/pkg/pretty"
)

type DeprecationMessage struct {
	Using		string
	Message		string
	Fix			string
	Note		string
}

var deprecationMap = map[string]DeprecationMessage{}

func base(id string, title, beforeUsing, afterUsing pretty.Segment, borderStyle pretty.BorderStyle, msg DeprecationMessage) {
	_, exists := deprecationMap[id]

	if exists {
		return
	}

	deprecationMap[id] = msg

	box := pretty.NewAutoBox()
	box.MinWidth = 60
	box.PaddingX = 2
	box.PaddingY = 1

	box.Border.Style = borderStyle

	box.AddBlock(pretty.Block{
		Align: pretty.AlignCenter,
		Segments: []pretty.Segment{
			title,
			pretty.InlineSegment{},
		},
	})

	if msg.Using != "" {
		box.AddBlock(pretty.Block{
			Align: pretty.AlignCenter,
			Segments: []pretty.Segment{
				beforeUsing,
				pretty.InlineSegment{},
				pretty.StyledTextBlockSegment{
					Raw: msg.Using,
				},
				pretty.InlineSegment{},
				afterUsing,
				pretty.InlineSegment{},
			},
		})
	}

	box.AddBlock(pretty.Block{
		Align: pretty.AlignCenter,
		Segments: []pretty.Segment{
			pretty.StyledTextBlockSegment{
				Raw: msg.Message,
			},
		},
	})

	if msg.Fix != "" {
		box.AddBlock(pretty.Block{
			Align: pretty.AlignCenter,
			Segments: []pretty.Segment{
				pretty.StyledTextBlockSegment{
					Raw: msg.Fix,
				},
			},
		})
	}

	if msg.Note != "" {
		box.AddBlock(pretty.Block{
			Align: pretty.AlignCenter,
			Segments: []pretty.Segment{
				pretty.StyledTextBlockSegment{
					Raw: msg.Note,
				},
			},
		})
	}

	fmt.Println(box.Render())
}

func Warn(id string, msg DeprecationMessage) {
	base(id,
		pretty.TextBlockSegment{
			Text: "🚧 Deprecation 🚧",
			Style: pretty.Style{
				Bold: true,
				Foreground: pretty.Basic(pretty.BrightYellow),
			},
		},
		pretty.TextBlockSegment{
			Text: "Please refrain from using",
		},
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
		pretty.BorderStyle{
			Color: pretty.Basic(pretty.BrightYellow),
		},
		msg,
	)
}

func Error(id string, msg DeprecationMessage) {
	base(id,
		pretty.TextBlockSegment{
			Text: "🚨 Breaking Change 🚨",
			Style: pretty.Style{
				Bold: true,
				Foreground: pretty.Basic(pretty.BrightRed),
			},
		},
		pretty.TextBlockSegment{
			Text: "Please stop using",
		},
		pretty.InlineSegment{
			Items: []pretty.Inline{
				pretty.Span{
					Text: "as it has been affected in a ",
				},
				pretty.Span{
					Text: "breaking change",
					Style: pretty.Style{
						Bold: true,
					},
				},
				pretty.Span{
					Text: ":",
				},
			},
		},
		pretty.BorderStyle{
			Color: pretty.Basic(pretty.BrightRed),
		},
		msg,
	)

	os.Exit(1)
}