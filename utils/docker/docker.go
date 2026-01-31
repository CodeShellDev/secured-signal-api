package docker

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/codeshelldev/gotl/pkg/docker"
	log "github.com/codeshelldev/gotl/pkg/logger"
)

func Init() {
	log.Info("Running ", os.Getenv("IMAGE_TAG"), " Image")
}

func Run(main func()) chan os.Signal {
	return docker.Run(main)
}

func Exit(code int) {
	log.Info("Exiting...")

	docker.Exit(code)
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
