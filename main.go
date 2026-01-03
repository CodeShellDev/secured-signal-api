package main

import (
	"os"
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

	logger.Info(`
	
	[1;34m┌────────────────────────────────────────────────┐[0m
	[1;34m│[0m [1;32m             🎄 Happy Holidays! 🎄            [0m [1;34m│[0m
	[1;34m│[0m                                                [1;34m│[0m
	[1;34m│[0m [0;37mThank you for using this project and for all  [0m [1;34m│[0m
	[1;34m│[0m [0;37mthe downloads, stars, and support this year.  [0m [1;34m│[0m
	[1;34m│[0m                                                [1;34m│[0m
	[1;34m│[0m [1;32mYour support truly means a lot — here's to    [0m [1;34m│[0m
	[1;34m│[0m [1;32man awesome year ahead! ✨                     [0m [1;34m│[0m
	[1;34m│[0m                                                [1;34m│[0m
	[1;34m│[0m [1;36m                 - CodeShell                  [0m [1;34m│[0m
	[1;34m└────────────────────────────────────────────────┘[0m
	`)

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
		port := config.SERVICE.PORT

		if strings.TrimSpace(port) != "" {
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
