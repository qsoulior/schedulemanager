package manager

import (
	"context"
	"log"

	"github.com/qsoulior/schedulemanager/internal/app"
	"github.com/qsoulior/schedulemanager/internal/repository"
	"github.com/qsoulior/schedulemanager/internal/usecase"
	"github.com/qsoulior/schedulemanager/pkg/mongodb"
	"github.com/qsoulior/schedulemanager/pkg/moodle"
)

func Run(config *app.Config) {
	mongo, err := mongodb.NewContext(config.Mongo.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Disconnect()
	mongoDatabase := mongo.Client.Database("app")

	moodleClient, err := moodle.NewClient(config.Moodle.Host, config.Moodle.Username, config.Moodle.Password)
	if err != nil {
		log.Fatal(err)
	}

	service := usecase.NewPlanService(repository.NewMongo(mongoDatabase), repository.NewMoodle(moodleClient))
	if err := service.AddSchedules(context.Background(), config.Moodle.CourseId); err != nil {
		log.Fatal(err)
	}
}
