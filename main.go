package main

import (
	"net/http"
	"os"

	config "github.com/codeshelldev/secured-signal-api/internals/config"
	reverseProxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	docker "github.com/codeshelldev/secured-signal-api/utils/docker"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var proxy reverseProxy.Proxy

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	log.Init(logLevel)

	docker.Init()

	config.Load()

	if config.DEFAULT.SERVICE.LOG_LEVEL != log.Level() {
		log.Init(config.DEFAULT.SERVICE.LOG_LEVEL)
	}

	log.Info("Initialized Logger with Level of ", log.Level())

	if log.Level() == "dev" {
		log.Dev("Welcome back Developer!")
		log.Dev("CTRL+S config to Print to Console")
	}

	config.Log()

	proxy = reverseProxy.Create(config.DEFAULT.API.URL)

	handler := proxy.Init()

	log.Info("Initialized Middlewares")

	addr := "0.0.0.0:" + config.DEFAULT.SERVICE.PORT

	log.Info("Server Listening on ", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	stop := docker.Run(func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", err.Error())
		}
	})

	<-stop

	docker.Shutdown(server)
}
