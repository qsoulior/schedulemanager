package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Context struct {
	Client *mongo.Client
}

func NewContext(uri string) (*Context, error) {
	c := new(Context)
	err := c.connect(uri)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Context) connect(uri string) error {
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

	c.Client = client
	return nil
}

func (c *Context) Disconnect() error {
	return c.Client.Disconnect(context.TODO())
}
