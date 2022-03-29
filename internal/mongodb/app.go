package mongodb

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppDB struct {
	client    *mongo.Client
	Schedules *SchedulesDriver
}

func (app *AppDB) Connect() error {
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

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	app.client = client
	db := client.Database("app")
	app.Schedules = NewSchedulesDriver(db)
	return nil
}

func (app *AppDB) Disconnect() error {
	return app.client.Disconnect(context.TODO())
}

func NewAppDB() (*AppDB, error) {
	appDriver := new(AppDB)
	err := appDriver.Connect()
	if err != nil {
		return nil, err
	}
	return appDriver, nil
}
