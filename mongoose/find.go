package mongoose

import (
	"context"
	"fmt"
	"time"

	"github.com/mrhid6/go-mongoose/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FindOne Searches one object and returns its value
func FindOne(filter bson.M, b interface{}) (err error) {
	// fmt.Println("Collection Name : ", utils.GetName(b))
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(utils.GetName(b))
	ctx, _ := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)

	res := collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return res.Err()
	}

	err = res.Decode(b)
	if err != nil {
		return err
	}

	return nil
}

// FindByID Searches by ID
func FindByID(id string, b interface{}) (err error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return FindOne(bson.M{
		"_id": objectID,
	}, b)
}

// FindByObjectID Searches by Object ID
func FindByObjectID(objectID bson.ObjectID, b interface{}) (err error) {
	return FindOne(bson.M{
		"_id": objectID,
	}, b)
}

// FindAllWithOptions Find all with options
func FindAllWithOptions(filter bson.M, option *options.FindOptionsBuilder, modelsOutArrayPtr interface{}) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(utils.GetName(modelsOutArrayPtr))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return err
	}
	err = cur.All(ctx, modelsOutArrayPtr)
	if err != nil {
		return err
	}
	return nil
}

// FindAll Get All Docs
func FindAll(filter bson.M, modelsOutArrayPtr interface{}) error {
	return FindAllWithOptions(filter, options.Find(), modelsOutArrayPtr)
}

// FindAllWithPagination Get All Docs with Pagination
func FindAllWithPagination(filter bson.M, start int64, count int64, modelsOutArrayPtr interface{}) error {

	opts := options.Find().SetSkip(start).SetLimit(count)

	return FindAllWithOptions(filter, opts, modelsOutArrayPtr)
}

func CountDocuments(collectionName string, filter bson.M) (int64, error) {
	mongo, err := Get()

	if err != nil {
		return 0, err
	}

	collection := mongo.Database.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count mods: %w", err)
	}

	return total, nil
}
