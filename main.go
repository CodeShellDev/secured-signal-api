package main

import (
	"net/http"
	"net/http/httputil"
	"os"

	proxy "github.com/codeshelldev/secured-signal-api/internals/proxy"
	middlewares "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares"
	config "github.com/codeshelldev/secured-signal-api/utils/config"
	docker "github.com/codeshelldev/secured-signal-api/utils/docker"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var proxy_last *httputil.ReverseProxy

var ENV *config.ENV_

func main() {
	logLevel := os.Getenv("LOG_LEVEL")

	log.Init(logLevel)

	config.Load()

	ENV = config.ENV

	if ENV.LOG_LEVEL != log.Level() {
		log.Init(ENV.LOG_LEVEL)
	}

	log.Info("Initialized Logger with Level of ", log.Level())

	proxy_last = proxy.Create(ENV.API_URL)

	body_m5 := middlewares.BodyMiddleware{
		Next: 	proxy_last,
	}

	temp_m4 := middlewares.TemplateMiddleware{
		Next: 	body_m5.Use(),
	}

	endp_m3 := middlewares.EndpointsMiddleware{
		Next: 	temp_m4.Use(),
	}

	auth_m2 := middlewares.AuthMiddleware{
		Next:   endp_m3.Use(),
	}

	serv_m1 := middlewares.ServeMiddleware{
		Next: auth_m2.Use(),
	}

	log_m0 := middlewares.LogMiddleware{
		Next: 	serv_m1.Use(),
	}

	log.Info("Initialized Proxy Handler")

	addr := "0.0.0.0:" + ENV.PORT

	log.Info("Server Listening on ", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: log_m0.Use(),
	}

	stop := docker.Run(func(){
		err := server.ListenAndServe()
		
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", err.Error())
		}
	})

	<-stop

	docker.Shutdown(server)
}