package prettylog

import (
	"fmt"
	"os"
	"strings"

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

func Deprecated(id string, msg DeprecationMessage) {
	if msg.Using == "" {
		msg.Using = decorateUsingPath("{b,i,bg=red}" + id + "{/}")
	}

	if msg.Note == "" {
		msg.Note = "\n{i}Update your config as {b}soon{/} as possible{/}"
	}

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

func Breaking(id string, msg DeprecationMessage) {
	if msg.Using == "" {
		msg.Using = decorateUsingPath("{b,i,bg=red}" + id + "{/}")
	}

	if msg.Note == "" {
		msg.Note = "\n{i}Update your config {b,fg=red}NOW!{/}{/}"
	}

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
					Text: "as it has been affected by a ",
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

func BreakingUsage(id string, msg DeprecationMessage) {
	if msg.Using == "" {
		msg.Using = decorateUsingPath("{b,i,bg=red}" + id + "{/}")
	}

	if msg.Note == "" {
		msg.Note = "\n{i}Update your config {b,fg=red}NOW!{/}{/}"
	}

	base(id,
		pretty.TextBlockSegment{
			Text: "🚨 Breaking Change 🚨",
			Style: pretty.Style{
				Bold: true,
				Foreground: pretty.Basic(pretty.BrightRed),
			},
		},
		pretty.TextBlockSegment{
			Text: "Please check your usage of",
		},
		pretty.InlineSegment{
			Items: []pretty.Inline{
				pretty.Span{
					Text: "as it has been affected by a ",
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

type StackOptions struct {
	Under 	int
	From 	int
	Count 	int
}

func GenericError(title string, err error) {
	GenericErrorWith(title, err, StackOptions{
		Under: 1,
		From: 3,
		Count: 6,
	})
}

func GenericErrorWith(title string, err error, opts StackOptions) {
	box := pretty.NewAutoBox()
	box.MinWidth = 60
	box.PaddingX = 2
	box.PaddingY = 1

	box.Border.Style = pretty.BorderStyle{
		Color: pretty.Basic(pretty.BrightRed),
	}

	box.AddBlock(pretty.Block{
		Align: pretty.AlignCenter,
		Segments: []pretty.Segment{
			pretty.StyledTextBlockSegment{
				Raw: title,
			},
			pretty.InlineSegment{},
		},
	})

	box.AddBlock(pretty.Block{
		Align: pretty.AlignCenter,
		Segments: []pretty.Segment{
			pretty.StyledTextBlockSegment{
				Raw: highlightKeywords(transformQuotes(err.Error(), func(open, close, inside string) string {
					return "{fg=green}" + open + inside + close + "{/}"
				})),
			},
		},
	})

	stack := prettyStack(opts.From, opts.Under, opts.Count)

	if stack != "" {
		box.AddBlock(pretty.Block{
			Align: pretty.AlignCenter,
			Segments: []pretty.Segment{
				pretty.InlineSegment{},
				pretty.StyledTextBlockSegment{
					Raw: stack,
				},
			},
		})
	}

	fmt.Println(box.Render())

	os.Exit(1)
}

func decorateUsingPath(path string) string {
	atRoot := !strings.Contains(path, ".")

	if atRoot {
		usingPrefix := "⇧ "
		usingSuffix := " (at root)"

		return usingPrefix + path + usingSuffix
	}

	return path
}