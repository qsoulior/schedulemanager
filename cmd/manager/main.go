package main

import (
	"os"

	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/app/manager"
)

func main() {
	log := app.NewLogger()

	configPath := os.Getenv("CONFIG_PATH")
	config, err := app.NewConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	manager.Run(config, log)
}
