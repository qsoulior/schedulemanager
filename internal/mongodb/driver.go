package mongodb

import (
	"context"

	"github.com/1asagne/schedulemanager/internal/schedule"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (driver *SchedulesDriver) GetOne(name string) (schedule.Schedule, error) {
	var schedule schedule.Schedule
	err := driver.collection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&schedule)
	return schedule, err
}

func (driver *SchedulesDriver) GetAll() ([]schedule.Schedule, error) {
	cursor, err := driver.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	schedules := make([]schedule.Schedule, 0)
	if err := cursor.All(context.TODO(), &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (driver *SchedulesDriver) DeleteAll() error {
	_, err := driver.collection.DeleteMany(context.TODO(), bson.D{})
	return err
}
