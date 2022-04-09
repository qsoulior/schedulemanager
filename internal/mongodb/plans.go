package mongodb

import (
	"context"
	"time"

	"github.com/1asagne/schedulemanager/internal/schedule"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Plan = schedule.Plan
type Schedule = schedule.Schedule

type PlansDriver struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewPlansDriver(db *mongo.Database) *PlansDriver {
	driver := new(PlansDriver)
	driver.db = db
	driver.collection = db.Collection("plans")
	return driver
}

func (driver *PlansDriver) AddSchedules(group string, schedule ...Schedule) error {
	updateOptions := options.Update().SetUpsert(true)
	_, err := driver.collection.UpdateOne(
		context.TODO(),
		bson.D{{"group", group}},
		bson.D{{"$set", bson.D{{"active", true}}}, {"$addToSet", bson.D{{"schedules", bson.D{{"$each", schedule}}}}}},
		updateOptions,
	)
	if err != nil {
		return err
	}
	return nil
}

func (driver *PlansDriver) GetSchedules(group string) ([]Schedule, error) {
	cursor, err := driver.collection.Aggregate(context.TODO(), mongo.Pipeline{
		bson.D{
			{"$match", bson.D{{"group", group}}},
		},
		bson.D{
			{"$unwind", "$schedules"},
		},
		bson.D{
			{"$replaceWith", "$schedules"},
		},
		bson.D{
			{"$sort", bson.D{{"modified", -1}}},
		},
	})
	if err != nil {
		return nil, err
	}
	schedules := make([]Schedule, 0)
	for cursor.Next(context.TODO()) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, err
}

func (driver *PlansDriver) GetScheduleLast(group string) (Schedule, error) {
	schedules, err := driver.GetSchedules(group)
	if err != nil {
		return Schedule{}, err
	}
	if len(schedules) > 0 {
		return schedules[0], nil
	}
	return Schedule{}, nil
}

type PlanInfo struct {
	Group    string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (driver *PlansDriver) GetInfo() ([]PlanInfo, error) {
	cursor, err := driver.collection.Aggregate(context.TODO(), mongo.Pipeline{
		bson.D{
			{"$unwind", "$schedules"},
		},
		bson.D{
			{"$group", bson.D{
				{"_id", "$group"},
				{"group", bson.D{{"$first", "$group"}}},
				{"modified", bson.D{{"$max", "$schedules.modified"}}},
			}},
		},
	})
	if err != nil {
		return nil, err
	}
	plansInfo := make([]PlanInfo, 0)
	for cursor.Next(context.TODO()) {
		var planInfo PlanInfo
		if err := cursor.Decode(&planInfo); err != nil {
			return nil, err
		}
		plansInfo = append(plansInfo, planInfo)
	}
	return plansInfo, nil
}
