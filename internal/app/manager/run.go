package manager

import (
	"context"

	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/repository"
	"github.com/qsoulior/schedulemanager/internal/usecase"
	"github.com/qsoulior/schedulemanager/pkg/mongodb"
	"github.com/qsoulior/schedulemanager/pkg/moodle"
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

	moodleClient, err := moodle.NewClient(config.Moodle.Host, config.Moodle.Username, config.Moodle.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Moodle client")
	}
	log.Info().Msg("Moodle client is created")

	service := usecase.NewPlanService(repository.NewMongo(mongoDatabase), repository.NewMoodle(moodleClient), log)
	if err := service.AddSchedules(context.Background(), config.Moodle.CourseId); err != nil {
		log.Fatal().Err(err).Msg("Failed to manage schedules")
	}
}
