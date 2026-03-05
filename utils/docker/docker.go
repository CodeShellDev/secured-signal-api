package docker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/codeshelldev/gotl/pkg/docker"
	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/pretty"
	httpserver "github.com/codeshelldev/gotl/pkg/server/http"
	"github.com/codeshelldev/secured-signal-api/utils/semver"
)

var VERSION *semver.Version

func Init() {
	imageTag := os.Getenv("IMAGE_TAG")

	if imageTag == "" {
		return
	}

	if semver.IsValid(imageTag) {
		v := semver.ParseSemver(imageTag)

		VERSION = &v

		logger.Info("Running ", VERSION.String(), " Image")

		if VERSION.Type != semver.FULL_RELEASE {
			box := pretty.NewAutoBox()

			box.Border.Style.Color = pretty.Basic(pretty.BrightBlue)
			box.Border.Chars = pretty.BorderChars{
				TopLeft: '+',
				BottomLeft: '+',
				TopRight: '+',
				BottomRight: '+',
				Horizontal: '─',
				Vertical: '│',
			}

			box.MinWidth = 60
			box.PaddingX = 2
			box.PaddingY = 1

			box.AddBlock(pretty.Block{
				Align: pretty.AlignCenter,
				Segments: []pretty.Segment{
					pretty.TextBlockSegment{
						Text: "🛠️  Pre-release Version 🛠️",
						Style: pretty.Style{
							Bold: true,
						},
					},
					pretty.InlineSegment{},
				},
			})

			box.AddBlock(pretty.Block{
				Align: pretty.AlignCenter,
				Segments: []pretty.Segment{
					pretty.StyledTextBlockSegment{
						Raw: "This is a" + 
							func() string { if VERSION.Type == semver.ALPHA_RELEASE { return "n" } else { return "" } }() + 
							" {i,b}" + string(VERSION.Type.Long()) + "{/}" + 
							func() string { if VERSION.Type != semver.RC_RELEASE { return " release" } else { return "" } }() + ", it may contain {b,fg=red}bugs{/} and ",
					},
					pretty.StyledTextBlockSegment{
						Raw: "some features may be {b,fg=bright_black}incomplete{/} or {b,fg=bright_yellow}unstable{/}",
					},
				},
			})

			box.AddBlock(pretty.Block{
				Align: pretty.AlignCenter,
				Segments: []pretty.Segment{
					pretty.InlineSegment{},
					pretty.StyledTextBlockSegment{
						Raw: "Encounter {b,fg=blue}issues{/}? Please {b,fg=blue}Report{/} them here:\n{b,u,fg=cyan}https://github.com/codeshelldev/secured-signal-api/issues{/}",
					},
				},
			})

			fmt.Println(box.Render())
		}
	} else {
		logger.Info("Running custom ", imageTag, " Image")
	}
}

func Run(main func()) chan os.Signal {
	return docker.Run(main)
}

func Exit(code int) {
	logger.Info("Exiting...")

	docker.Exit(code)
}

func Shutdown(server *httpserver.HttpServer) {
	logger.Info("Shutdown signal received")

	logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		logger.Fatal("Server shutdown failed: ", err.Error())

		logger.Info("Server exited forcefully")
	} else {
		logger.Info("Server exited gracefully")
	}
}
