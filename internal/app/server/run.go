package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/controller/http"
	"github.com/qsoulior/schedulemanager/internal/repository"
	"github.com/qsoulior/schedulemanager/internal/usecase"
	"github.com/qsoulior/schedulemanager/pkg/mongodb"
)

func Run(config *app.Config) {
	mongo, err := mongodb.NewContext(config.Mongo.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Disconnect()
	mongoDatabase := mongo.Client.Database("app")

	service := usecase.NewPlanService(repository.NewMongo(mongoDatabase), nil)
	controller := http.NewPlanController(service)

	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:        500,
		Expiration: 1 * time.Minute,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Server.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type",
		MaxAge:       300,
	}))

	api := app.Group("/api")
	api.Get("/schedules", controller.GetSchedules)
	api.Get("/info", controller.GetPlansInfo)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.Server.Port)))
}
