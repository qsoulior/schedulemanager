package mongodb

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/1asagne/schedulemanager/internal/schedule"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SchedulesDriver struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewSchedulesDriver(db *mongo.Database) *SchedulesDriver {
	driver := new(SchedulesDriver)
	driver.db = db
	driver.collection = db.Collection("schedules")
	return driver
}

func (driver *SchedulesDriver) InsertOne(file []schedule.File) error {
	_, err := driver.collection.InsertOne(context.TODO(), file)
	return err
}

func (driver *SchedulesDriver) InsertMany(schedules []schedule.Schedule) error {
	documents := make([]interface{}, len(schedules))
	for i := range schedules {
		documents[i] = schedules[i]
	}
	_, err := driver.collection.InsertMany(context.TODO(), documents)
	return err
}

func (driver *SchedulesDriver) DeleteAll() error {
	_, err := driver.collection.DeleteMany(context.TODO(), bson.D{})
	return err
}

func SaveFiles(files []schedule.Schedule) error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("MONGODB_URI is missing in enviroment")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	database := client.Database("app")

	schedules := NewSchedulesDriver(database)
	if err := schedules.DeleteAll(); err != nil {
		return err
	}
	if err := schedules.InsertMany(files); err != nil {
		return err
	}
	return nil
}
