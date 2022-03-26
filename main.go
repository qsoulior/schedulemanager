package main

import (
	"log"

	"github.com/1asagne/schedulemanager/internal/mongodb"
	"github.com/1asagne/schedulemanager/internal/moodle"
	"github.com/1asagne/schedulemanager/internal/schedule"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatal(err)
		return
	}

	scheduleFiles, err := moodle.DownloadFiles()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("Getting files completed\n")

	scheduleFilesParsed, err := schedule.ParseFiles(scheduleFiles)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("Parsing files completed\n")

	if err := mongodb.SaveFiles(scheduleFilesParsed); err != nil {
		log.Fatal(err)
		return
	}
	log.Print("DB client initialization completed\n")
}
