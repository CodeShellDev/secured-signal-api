package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
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

		server.ListenAndServer()
	})

	<-stop

	db.Close()
	docker.Shutdown(server)
}
