package config

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI  string
	Name string
}

func SetupMongoDataBase(ctx context.Context) (*mongo.Database, error) {
	config := &MongoConfig{}

	err := viper.UnmarshalKey("mongo.database", config)

	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	return client.Database(config.Name), nil

}
