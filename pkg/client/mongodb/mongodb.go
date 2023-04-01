package mongodb

import (
	"context"
	"fmt"
	"github.com/Salavei/golang_websockets/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, conf config.StorageConfig) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if conf.MongoDB.Username == "" && conf.MongoDB.Password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", conf.MongoDB.Host, conf.MongoDB.Port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.MongoDB.Username, conf.MongoDB.Password,
			conf.MongoDB.Host, conf.MongoDB.Port)
	}
	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if conf.MongoDB.AuthDB == "" {
			conf.MongoDB.AuthDB = conf.MongoDB.Database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: conf.MongoDB.AuthDB,
			Username:   conf.MongoDB.Username,
			Password:   conf.MongoDB.Password,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB due to error: %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB due to error: %v", err)
	}

	return client.Database(conf.MongoDB.Database), nil
}
