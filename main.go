package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/pretty"
	httpserver "github.com/codeshelldev/gotl/pkg/server/http"
	config "github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/db"
	reverseProxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	"github.com/codeshelldev/secured-signal-api/internals/scheduler"
	docker "github.com/codeshelldev/secured-signal-api/utils/docker"
	"github.com/codeshelldev/secured-signal-api/utils/logging"
	"github.com/codeshelldev/secured-signal-api/utils/prettylog"
)

var proxy reverseProxy.Proxy

func catchPanic(fn func()) {
	defer func() {
		r := recover()

		if r != nil {
			switch v := r.(type) {
			case error:
			default:
				prettylog.GenericErrorWith("{b,fg=red}🚨 PANIC 🚨{/}", errors.New("encountered a panic:\n\n" + fmt.Sprint(v)), prettylog.StackOptions{
					Count: 8,
					Under: 2,
					From: 3,
				})
			}
		}
	}()

	fn()
}

func main() {
	catchPanic(m)
}

func m() {
	logging.Init(os.Getenv("LOG_LEVEL"))

	docker.Init()

	titleOfNextMessage := "Happy Easter!"

	// c'mon this is way too late, let's give 2.4 a pass too
	if time.Now().Format("2.1") == "1.4" || time.Now().Format("2.1") == "2.4" {
		// TODO remove april fools
		box := pretty.NewAutoBox()
		box.MinWidth = 60
		box.PaddingX = 2
		box.PaddingY = 1

		box.Border.Style = pretty.BorderStyle{
			Color: pretty.Basic(pretty.Blue),
		}
		box.AddBlock(pretty.Block{
			Align: pretty.AlignCenter,
			Segments: []pretty.Segment{
				pretty.StyledTextBlockSegment{
					Raw: `{b}📣 Introducing: {fg=magenta}SSA AI{/} 📣{/}

We are excited to introduce to you our newest {b}solution{/}:

✨ {b,fg=blue}Secured Signal API AI{/} ✨

{b,fg=blue}SSA AI{/} extends the Secured Signal API with
message handling, intent parsing, and full-stack signal participation

It processes every request through a continuous {b}identity-linked{/} inference mesh 
that ensures {b}no message is ever truly "unobserved" 👀{/}

Furthermore, our backends are now {b}99% AI{/} and have been rewritten into {b}JavaScript{/} 🎉
To further emphasize AI we have also decided to {b}let AI handle authentication{/} 🔐`,
				},
			},
		})

		titleOfNextMessage = "APRIL FOOLS – Happy Easter!"

		fmt.Println(box.Render())

		time.Sleep(3 * time.Second)
	}

	// TODO remove greeting
	box := pretty.NewAutoBox()
	box.MinWidth = 60
	box.PaddingX = 2
	box.PaddingY = 1

	box.Border.Style = pretty.BorderStyle{
		Color: pretty.Basic(pretty.Blue),
	}
	box.AddBlock(pretty.Block{
		Align: pretty.AlignCenter,
		Segments: []pretty.Segment{
			pretty.StyledTextBlockSegment{
				Raw: `{b,fg=blue}🐰 ` + titleOfNextMessage + ` 🐰{/}

…and {b,fg=red}thank you{/} ❤️  for using {b}Secured Signal API{/}.

Since the last {i}2 months{/} {i}(wow, has been it been long…){/}

we have started to gain a lot of {b,fg=yellow}pulls{/} and {b,fg=yellow}starts{/} ⭐️.
We even got some {b,fg=blue}issues{/} opened by you all 🥳!

– CodeShell 🐢`,
			},
		},
	})

	fmt.Println(box.Render())

	config.Load()

	config.Validate()

	if config.DEFAULT.SERVICE.LOG_LEVEL != logger.Level() {
		logging.Init(config.DEFAULT.SERVICE.LOG_LEVEL)
	}

	logger.Info("Initialized Logger with Level of ", logger.Level())

	logging.Setup()

	if logger.Level() == "dev" {
		logger.Dev("Welcome back, Developer!")
		logger.Dev("Resave config to trigger reload and print")
	}

	config.Log()

	db.Init()

	scheduler.Start()

	proxy = reverseProxy.Create((*url.URL)(config.DEFAULT.API.URL))

	handler := proxy.Init()

	logger.Info("Initialized Middlewares")

	ports := []string{}

	for _, config := range config.ENV.CONFIGS {
		port := strings.TrimSpace(config.SERVICE.PORT)

		if port != "" && !slices.Contains(ports, port) {
			ports = append(ports, port)
		}
	}

	server := httpserver.Create(handler, "0.0.0.0", ports...)

	server.ErrorLog = logger.StdError()
	server.InfoLog = logger.StdInfo()

	stop := docker.Run(func() {
		if logger.IsDebug() && len(ports) > 1 {
			logger.Debug("Server started with ", len(ports), " listeners on ", httpserver.PortsToRangeString(ports))
		} else {
			logger.Info("Server listening on ", httpserver.PortsToRangeString(ports))
		}

		catchPanic(server.ListenAndServer)
	})

	<-stop

	db.Close()
	docker.Shutdown(server)
}
