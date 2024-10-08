package mongoose

import (
	"context"
	"time"

	"github.com/mrhid6/go-mongoose/mutility"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateByID Updates by ID
func UpdateByID(model interface{}) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	_, err = collection.ReplaceOne(ctx, bson.M{
		"_id": mutility.GetID(model),
	}, model)

	if err != nil {
		return err
	}
	return nil
}

func UpdateDataByID(model interface{}, update interface{}) error {

	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": mutility.GetID(model),
	}, update)

	if err != nil {
		return err
	}
	return nil
}
