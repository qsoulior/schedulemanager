package mongodb

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	client *mongo.Client
	Plans  *PlansDriver
}

func (app *App) Connect() error {
	uri := os.Getenv("MONGODB_CONNSTRING")
	if uri == "" {
		return errors.New("MONGODB_CONNSTRING is missing in enviroment")
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
	app.Plans = NewPlansDriver(db)
	return nil
}

func (app *App) Disconnect() error {
	return app.client.Disconnect(context.TODO())
}

func NewApp() (*App, error) {
	appDriver := new(App)
	err := appDriver.Connect()
	if err != nil {
		return nil, err
	}
	return appDriver, nil
}
