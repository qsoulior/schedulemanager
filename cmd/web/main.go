package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/qsoulior/schedulemanager/internal/mongodb"
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

	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:        500,
		Expiration: 1 * time.Minute,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("API_ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type",
		MaxAge:       300,
	}))

	api := app.Group("/api")

	schedulesHandler := func(c *fiber.Ctx) error {
		group := c.Query("group")

		handleError := func(err error) *fiber.Error {
			if err == mongodb.ErrNoSchedules {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			errorLog.Println(err)
			return fiber.ErrInternalServerError
		}

		if c.Query("newest") == "" {
			if schedules, err := db.Plans.GetSchedules(group); err != nil {
				return handleError(err)
			} else {
				return c.JSON(schedules)
			}
		}

		if schedule, err := db.Plans.GetScheduleLast(group); err != nil {
			return handleError(err)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	errorLog.Fatal(app.Listen(":" + port))
}
