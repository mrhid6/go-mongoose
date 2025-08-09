package mongoose

import (
	"context"
	"time"

	"github.com/mrhid6/go-mongoose/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateByID Updates by ID
func UpdateByID(model interface{}) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(utils.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	_, err = collection.ReplaceOne(ctx, bson.M{
		"_id": utils.GetID(model),
	}, model)

	if err != nil {
		return err
	}
	return nil
}

// Deprecated: Use UpdateModelData instead.
func UpdateDataByID(model interface{}, update interface{}) error {

	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(utils.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": utils.GetID(model),
	}, update)

	if err != nil {
		return err
	}
	return nil
}

func UpdateModelData(model interface{}, updateData bson.M) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(utils.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": utils.GetID(model),
	}, bson.M{"$set": updateData})

	if err != nil {
		return err
	}
	return nil
}
