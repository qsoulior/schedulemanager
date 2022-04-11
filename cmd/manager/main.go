package main

import (
	"log"
	"os"

	"github.com/1asagne/schedulemanager/internal/mongodb"
	"github.com/1asagne/schedulemanager/internal/moodle"
	"github.com/1asagne/schedulemanager/internal/schedule"
	"github.com/joho/godotenv"
)

var infoLog, errorLog *log.Logger

func init() {
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load("dev.env")
	if err != nil {
		errorLog.Fatal(err)
	}
}

func main() {
	db, err := mongodb.NewApp()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Disconnect()
	infoLog.Print("DB initialization completed\n")

	scheduleFiles, err := moodle.DownloadFiles(db)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Printf("Schedules downloading completed. New schedules: %d.\n", len(scheduleFiles))

	if len(scheduleFiles) > 0 {
		plans, err := schedule.ParseFiles(scheduleFiles)
		if err != nil {
			errorLog.Fatal(err)
		}
		infoLog.Print("Schedules parsing completed\n")
		for _, plan := range plans {
			if err := db.Plans.AddSchedules(plan.Group, plan.Schedules...); err != nil {
				errorLog.Fatal(err)
			}
		}
		infoLog.Print("Parsed schedules saving completed\n")
	}
}
