package docker

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
)

var stop chan os.Signal

func Init() {
	log.Info("Running ", os.Getenv("IMAGE_TAG"), " Image")
}

func Run(main func()) chan os.Signal {
	stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go main()

	return stop
}

func Exit(code int) {
	log.Info("Exiting...")

	os.Exit(code)

	stop <- syscall.SIGTERM
}

func Shutdown(server *http.Server) {
	log.Info("Shutdown signal received")

	log.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Fatal("Server shutdown failed: ", err.Error())

		log.Info("Server exited forcefully")
	} else {
		log.Info("Server exited gracefully")
	}
}
