package main

import (
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
	"github.com/codeshelldev/secured-signal-api/utils/stdlog"
)

var proxy reverseProxy.Proxy

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	logger.Init(logLevel)

	docker.Init()

	config.Load()

	if config.DEFAULT.SERVICE.LOG_LEVEL != logger.Level() {
		logger.Init(config.DEFAULT.SERVICE.LOG_LEVEL)
	}

	logger.Info("Initialized Logger with Level of ", logger.Level())

	if logger.Level() == "dev" {
		logger.Dev("Welcome back, Developer!")
		logger.Dev("CTRL+S config to Print to Console")
	}

	config.Log()

	db.Init()

	scheduler.Start()

	proxy = reverseProxy.Create(config.DEFAULT.API.URL.URL)

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

	server.ErrorLog = stdlog.ErrorLog
	server.InfoLog = stdlog.DebugLog

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
