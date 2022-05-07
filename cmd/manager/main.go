package main

import (
	"log"
	"os"

	"github.com/qsoulior/schedulemanager/internal/mongodb"
	"github.com/qsoulior/schedulemanager/internal/moodle"
	"github.com/qsoulior/schedulemanager/internal/schedule"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

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
