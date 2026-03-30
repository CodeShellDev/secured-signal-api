package docker

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

	if imageTag != "" {
		if semver.IsValid(imageTag) {
			v := semver.ParseSemver(imageTag)

			VERSION = &v

			logger.Info("Running ", VERSION.String(), " Image")
			logger.Debug("Image built at ", )

			if VERSION.Type != semver.FULL_RELEASE {
				box := pretty.NewAutoBox()

				box.Border.Style.Color = pretty.Basic(pretty.BrightBlue)

				box.MinWidth = 60
				box.PaddingX = 2
				box.PaddingY = 1

				box.AddBlock(pretty.Block{
					Align: pretty.AlignCenter,
					Segments: []pretty.Segment{
						pretty.TextBlockSegment{
							Text: "🔬 Pre-Release 🔬",
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
							Raw: "Encounter {b,fg=blue}issues{/}? Please {b,fg=blue}Report{/} them here:\n\n{b,u,fg=cyan}https://codeshelldev.github.io/secured-signal-api/bug?v=" + VERSION.String() + "{/}",
						},
					},
				})

				fmt.Println(box.Render())
			}
		} else {
			logger.Info("Running custom ", imageTag, " Image")
		}
	}
	
	buildTimestamp := os.Getenv("BUILD_TIME")
	buildAt, err := strconv.Atoi(buildTimestamp)

	if err == nil {
		buildTime := time.Unix(int64(buildAt), 0)

		if !buildTime.After(time.Now()) {
			logger.Info("Built time: ", buildTime.Local().Format("01.02.06 15:04:05"))
		}
	}

	commit := os.Getenv("GIT_COMMIT")

	if commit != "" {
		re, err := regexp.Compile(`\b[0-9a-f]{40}\b`)

		if err == nil && re.MatchString(commit) {
			logger.Info("Commit: ", commit)
		}
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
