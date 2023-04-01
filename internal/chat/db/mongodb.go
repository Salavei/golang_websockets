package db

import (
	"context"
	"fmt"
	"github.com/Salavei/golang_websockets/internal/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type db struct {
	collection *mongo.Collection
}

func (d *db) SendMessage(ctx context.Context, user string, msg chat.Message) (string, error) {

	log.Println("create message")
	_, err := d.collection.InsertOne(ctx, msg)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("failed to create message due to error: %v", err)
	}
	return "", nil
}

func (d *db) ShowMessage(ctx context.Context) (msgs []chat.Message, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var msg chat.Message
		err = cursor.Decode(&msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return msgs, nil
}

func NewStorage(database *mongo.Database, collection string) chat.Storage {

	return &db{
		collection: database.Collection(collection),
	}
}
