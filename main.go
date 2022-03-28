package main

import (
	"log"

	"github.com/1asagne/schedulemanager/internal/mongodb"
	"github.com/1asagne/schedulemanager/internal/moodle"
	"github.com/1asagne/schedulemanager/internal/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatal(err)
	}

	scheduleFiles, err := moodle.DownloadFiles()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Schedules downloading completed\n")

	scheduleFilesParsed, err := schedule.ParseFiles(scheduleFiles)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Schedules parsing completed\n")

	dbAppInstance, err := mongodb.NewAppInstance()
	if err != nil {
		log.Fatal(err)
	}
	defer dbAppInstance.Disconnect()
	log.Print("DB initialization completed\n")

	if err := dbAppInstance.SaveFiles(scheduleFilesParsed); err != nil {
		log.Fatal(err)
	}
	log.Print("Schedules saving completed\n")

	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Test")
	})

	webApp.Listen(":3000")
}
