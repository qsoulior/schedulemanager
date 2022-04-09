package main

import (
	"log"
	"os"
	"time"

	"github.com/1asagne/schedulemanager/internal/mongodb"
	"github.com/1asagne/schedulemanager/internal/moodle"
	"github.com/1asagne/schedulemanager/internal/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:        500,
		Expiration: 1 * time.Minute,
	}))

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		errorLog.Fatal("API_TOKEN is missing in enviroment")
	}

	api := app.Group("/api", func(c *fiber.Ctx) error {
		if c.Query("token") == apiToken {
			return c.Next()
		}
		infoLog.Printf("Unauthorized request from %s\n", c.IP())
		return fiber.ErrUnauthorized
	})

	schedulesHandler := func(c *fiber.Ctx) error {
		group := c.Query("group")
		if c.Query("last") == "" {
			if schedules, err := db.Plans.GetSchedules(group); err != nil {
				errorLog.Println(err)
				return fiber.ErrInternalServerError
			} else {
				return c.JSON(schedules)
			}
		}
		if schedule, err := db.Plans.GetScheduleLast(group); err != nil {
			errorLog.Println(err)
			return fiber.ErrInternalServerError
		} else {
			return c.JSON(schedule)
		}
	}
	infoHandler := func(c *fiber.Ctx) error {
		info, err := db.Plans.GetInfo()
		if err != nil {
			errorLog.Println(err)
			return fiber.ErrInternalServerError
		}
		return c.JSON(info)
	}

	api.Get("/schedules", schedulesHandler)
	api.Get("/info", infoHandler)

	errorLog.Fatal(app.Listen(":3000"))
}
