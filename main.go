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

	web := fiber.New()
	web.Get("/schedules", func(c *fiber.Ctx) error {
		group := c.Query("group")
		if c.Query("last") == "" {
			if schedules, err := db.Plans.GetSchedules(group); err != nil {
				errorLog.Println(err)
				return c.SendStatus(500)
			} else {
				return c.JSON(schedules)
			}
		}
		if schedule, err := db.Plans.GetScheduleLast(group); err != nil {
			errorLog.Println(err)
			return c.SendStatus(500)
		} else {
			return c.JSON(schedule)
		}
	})
	web.Get("/info", func(c *fiber.Ctx) error {
		info, err := db.Plans.GetInfo()
		if err != nil {
			errorLog.Println(err)
			return c.SendStatus(500)
		}
		return c.JSON(info)
	})

	web.Listen(":3000")
}
