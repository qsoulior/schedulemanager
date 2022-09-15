package repository

import (
	"context"

	"github.com/qsoulior/schedulemanager/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	collection *mongo.Collection
}

func NewMongo(db *mongo.Database) PlanDatabase {
	collection := db.Collection("plans")
	return &Mongo{collection}
}

func (repo *Mongo) AddSchedules(ctx context.Context, group string, schedules ...entity.Schedule) error {
	updateOptions := options.Update().SetUpsert(true)
	_, err := repo.collection.UpdateOne(
		ctx,
		bson.D{{"group", group}},
		bson.D{{"$set", bson.D{{"active", true}}}, {"$addToSet", bson.D{{"schedules", bson.D{{"$each", schedules}}}}}},
		updateOptions,
	)
	return err
}

func (repo *Mongo) GetSchedules(ctx context.Context, group string) ([]entity.Schedule, error) {
	cursor, err := repo.collection.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{{"group", group}, {"active", true}}},
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
	schedules := make([]entity.Schedule, 0)
	for cursor.Next(ctx) {
		schedule := new(entity.Schedule)
		if err := cursor.Decode(schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, *schedule)
	}
	return schedules, err
}

func (repo *Mongo) GetLatestSchedule(ctx context.Context, group string) (*entity.Schedule, error) {
	schedules, err := repo.GetSchedules(ctx, group)
	if err != nil {
		return nil, err
	}
	return &schedules[0], nil
}

func (repo *Mongo) GetPlansInfo(ctx context.Context) ([]entity.PlanInfo, error) {
	cursor, err := repo.collection.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{{"active", true}}},
		},
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
	plansInfo := make([]entity.PlanInfo, 0)
	for cursor.Next(ctx) {
		planInfo := new(entity.PlanInfo)
		if err := cursor.Decode(planInfo); err != nil {
			return nil, err
		}
		plansInfo = append(plansInfo, *planInfo)
	}
	return plansInfo, nil
}

func (repo *Mongo) DeactivatePlan(ctx context.Context, group string) error {
	_, err := repo.collection.UpdateOne(
		ctx,
		bson.D{{"group", group}},
		bson.D{{"$set", bson.D{{"active", false}}}},
	)
	return err
}
