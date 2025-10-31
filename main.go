package main

import (
	"net/http"
	"os"

	config "github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	reverseProxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	docker "github.com/codeshelldev/secured-signal-api/utils/docker"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var proxy reverseProxy.Proxy

var ENV *structure.ENV

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	log.Init(logLevel)

	docker.Init()

	config.Load()

	ENV = config.ENV

	if ENV.LOG_LEVEL != log.Level() {
		log.Init(ENV.LOG_LEVEL)
	}

	log.Info("Initialized Logger with Level of ", log.Level())

	proxy = reverseProxy.Create(ENV.API_URL)

	handler := proxy.Init()

	log.Info("Initialized Middlewares")

	addr := "0.0.0.0:" + ENV.PORT

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
