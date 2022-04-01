package mongodb

import (
	"context"
	"time"

	"github.com/1asagne/schedulemanager/internal/schedule"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SchedulesDriver struct {
	db         *mongo.Database
	collection *mongo.Collection
}

type ScheduleInfo struct {
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
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
	err := driver.collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "name", Value: name}}).Decode(&schedule)
	return schedule, err
}

func (driver *SchedulesDriver) GetAll() ([]schedule.Schedule, error) {
	cursor, err := driver.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	results := make([]schedule.Schedule, 0)

	for cursor.Next(context.TODO()) {
		var result schedule.Schedule
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	defer cursor.Close(context.TODO())
	return results, nil
}

func (driver *SchedulesDriver) GetAllInfo() ([]ScheduleInfo, error) {
	opts := options.Find()
	opts.SetProjection(bson.M{"name": true, "modified": true, "_id": false})
	cursor, err := driver.collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	results := make([]ScheduleInfo, 0)
	for cursor.Next(context.TODO()) {
		var result ScheduleInfo
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	defer cursor.Close(context.TODO())
	return results, nil
}

func (driver *SchedulesDriver) DeleteAll() error {
	_, err := driver.collection.DeleteMany(context.TODO(), bson.D{})
	return err
}
