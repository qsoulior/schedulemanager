package mongodb

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/1asagne/schedulemanager/internal/schedule"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppInstance struct {
	db     *mongo.Database
	client *mongo.Client
}

func (app *AppInstance) Connect() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("MONGODB_URI is missing in enviroment")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return err
	}
	app.client = client
	app.db = client.Database("app")
	return nil
}

func (app *AppInstance) Disconnect() error {
	return app.client.Disconnect(context.TODO())
}

func NewAppInstance() (*AppInstance, error) {
	appDriver := new(AppInstance)
	err := appDriver.Connect()
	if err != nil {
		return nil, err
	}
	return appDriver, nil
}

func (app *AppInstance) SaveFiles(files []schedule.Schedule) error {
	schedules := NewSchedulesDriver(app.db)
	if err := schedules.DeleteAll(); err != nil {
		return err
	}
	if err := schedules.InsertMany(files); err != nil {
		return err
	}
	return nil
}

func (app *AppInstance) GetFiles() ([]schedule.Schedule, error) {
	schedules := NewSchedulesDriver(app.db)
	return schedules.GetAll()
}

func (app *AppInstance) GetFile(name string) (schedule.Schedule, error) {
	schedules := NewSchedulesDriver(app.db)
	return schedules.GetOne(name)
}
