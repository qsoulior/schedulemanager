package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/controller/http"
	"github.com/qsoulior/schedulemanager/internal/repository"
	"github.com/qsoulior/schedulemanager/internal/usecase"
	"github.com/qsoulior/schedulemanager/pkg/mongodb"
	"github.com/rs/zerolog"
)

func Run(config *app.Config, log *zerolog.Logger) {
	mongo, err := mongodb.NewContext(config.Mongo.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection failed")
	}
	log.Info().Msg("Connected to database")
	defer mongo.Disconnect()
	mongoDatabase := mongo.Client.Database("app")

	service := usecase.NewPlanService(repository.NewMongo(mongoDatabase), nil, log)
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

	log.Info().Str("port", strconv.Itoa(config.Server.Port)).Msg("Starting http server")
	if err := app.Listen(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal().Err(err).Msg("Server listening failed")
	}
}
