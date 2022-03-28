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

	scheduleFiles, err := moodle.DownloadFiles()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Schedules downloading completed\n")

	scheduleFilesParsed, err := schedule.ParseFiles(scheduleFiles)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Schedules parsing completed\n")

	dbAppInstance, err := mongodb.NewAppInstance()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbAppInstance.Disconnect()
	infoLog.Print("DB initialization completed\n")

	if err := dbAppInstance.SaveFiles(scheduleFilesParsed); err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Schedules saving completed\n")

	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		names, err := dbAppInstance.GetNames()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.JSON(names)
	})

	webApp.Listen(":3000")
}
