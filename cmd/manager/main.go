package main

import (
	"flag"
	"log"

	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/app/manager"
)

func main() {
	configPath := flag.String("c", "", "configuration file path")
	flag.Parse()

	if *configPath == "" {
		flag.PrintDefaults()
		return
	}
	config, err := app.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	manager.Run(config)
}
