package database

import (
	"context"
	"home/jonganebski/github/fibersteps-server/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var env = config.Config("APP_ENV")
var mongoURI = config.Config("MONGO_URI_" + env)
var dbName = config.Config("DB_NAME")

// Mongo is mongo instance variable
var Mongo mongoInstance

// Connect connects to MongoDB
func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	Mongo = mongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}
