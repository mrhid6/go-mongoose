package mongoose

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/mrhid6/go-mongoose/mutility"

	"go.mongodb.org/mongo-driver/mongo"
)

// InsertOne This will insert just one Data
func InsertOne(modelPtr interface{}) (res *mongo.InsertOneResult, err error) {

	if mutility.IsPointer(modelPtr) {
		return nil, errors.New("model should be a Pointer")
	}

	mongo, err := Get()

	if err != nil {
		return nil, err
	}

	collection := mongo.Database.Collection(mutility.GetName(modelPtr))
	ctx, _ := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)

	res, err = collection.InsertOne(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(modelPtr).Elem().FieldByName("ID")
	val.Set(reflect.ValueOf(res.InsertedID))
	return res, err
}

// InsertMany This will insert multiple Data
// TODO Find a way to pass pointer and attach its ID to the respective array elements
func InsertMany(models []interface{}) (res *mongo.InsertManyResult, err error) {
	if len(models) == 0 {
		return nil, errors.New("the length of Model Array is 0")
	}

	if mutility.IsPointer(models) {
		return nil, errors.New("models should be a Pointer")
	}

	mongo, err := Get()

	if err != nil {
		return nil, err
	}

	collection := mongo.Database.Collection(mutility.GetName(models))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	// iM := make([]interface{}, 0)
	// iM = append(iM, models)
	res, err = collection.InsertMany(ctx, models)
	if err != nil {
		return nil, err
	}
	return res, err
}
