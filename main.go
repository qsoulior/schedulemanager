package main

import (
	"log"
	"os"

	"github.com/1asagne/schedulemanager/internal/mongodb"
	"github.com/1asagne/schedulemanager/internal/moodle"
	"github.com/1asagne/schedulemanager/internal/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load("dev.env")
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := mongodb.NewAppDB()
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
		scheduleFilesParsed, err := schedule.ParseFiles(scheduleFiles)
		if err != nil {
			errorLog.Fatal(err)
		}
		infoLog.Print("Schedules parsing completed\n")

		if err := db.Schedules.InsertMany(scheduleFilesParsed); err != nil {
			errorLog.Fatal(err)
		}
		infoLog.Print("Parsed schedules saving completed\n")
	}

	web := fiber.New()
	web.Get("/all", func(c *fiber.Ctx) error {
		files, err := db.Schedules.GetAll()
		if err != nil {
			return c.SendStatus(500)
		}
		return c.JSON(files)
	})
	web.Get("/info", func(c *fiber.Ctx) error {
		names, err := db.Schedules.GetAllInfo()
		if err != nil {
			return c.SendStatus(500)
		}
		return c.JSON(names)
	})

	web.Listen(":3000")
}
