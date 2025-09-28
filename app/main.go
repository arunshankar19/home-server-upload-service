package main

import (
	"log"

	"github.com/arunshankar19/home-server-common-utils/logger"
	"github.com/arunshankar19/home-server-upload-service/internal/config"
)

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Println("failed to initialize logger with error", err)
		panic(err)
	}

	appConfig, err := config.NewAppConfig()
	if err != nil {
		l.Error("failed to get all required envs", map[string]any{
			"error": err,
		})
		panic(err)
	}

	l.Info("application configs read successfully", map[string]any{
		"appConfig": appConfig,
	})
}
