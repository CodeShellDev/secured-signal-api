package main

import (
	"os"
	"slices"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	config "github.com/codeshelldev/secured-signal-api/internals/config"
	reverseProxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	httpServer "github.com/codeshelldev/secured-signal-api/internals/server"
	docker "github.com/codeshelldev/secured-signal-api/utils/docker"
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
		logger.Dev("Welcome back Developer!")
		logger.Dev("CTRL+S config to Print to Console")
	}

	config.Log()

	proxy = reverseProxy.Create(config.DEFAULT.API.URL)

	handler := proxy.Init()

	logger.Info("Initialized Middlewares")

	ports := []string{}

	for _, config := range config.ENV.CONFIGS {
		port := strings.TrimSpace(config.SERVICE.PORT)

		if port != "" && !slices.Contains(ports, port) {
			ports = append(ports, port)
		}
	}

	server := httpServer.Create(handler, "0.0.0.0", ports...)

	stop := docker.Run(func() {
		if logger.IsDebug() && len(ports) > 1 {
			logger.Debug("Server started with ", len(ports), " listeners on ", httpServer.PortsToRangeString(ports))
		} else {
			logger.Info("Server listening on ", httpServer.PortsToRangeString(ports))
		}

		server.ListenAndServer()
	})

	<-stop

	docker.Shutdown(server)
}
