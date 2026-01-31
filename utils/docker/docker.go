package docker

import (
	"context"
	"os"
	"time"

	"github.com/codeshelldev/gotl/pkg/docker"
	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/server"
)

func Init() {
	logger.Info("Running ", os.Getenv("IMAGE_TAG"), " Image")
}

func Run(main func()) chan os.Signal {
	return docker.Run(main)
}

func Exit(code int) {
	logger.Info("Exiting...")

	docker.Exit(code)
}

func Shutdown(server *server.Server) {
	logger.Info("Shutdown signal received")

	logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		logger.Fatal("Server shutdown failed: ", err.Error())

		logger.Info("Server exited forcefully")
	} else {
		logger.Info("Server exited gracefully")
	}
}
