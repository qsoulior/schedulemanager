package main

import (
	"flag"

	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/app/manager"
)

func main() {
	log := app.NewLogger()

	configPath := flag.String("c", "", "configuration file path")
	flag.Parse()

	if *configPath == "" {
		flag.PrintDefaults()
		return
	}
	config, err := app.NewConfig(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	manager.Run(config, log)
}
