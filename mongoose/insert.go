package mongoose

import (
	"context"
	"errors"
	"reflect"
	"time"

	"go-mongoose/mutility"

	"go.mongodb.org/mongo-driver/mongo"
)

//InsertOne This will insert just one Data
func InsertOne(modelPtr interface{}) (res *mongo.InsertOneResult, err error) {
	collection := Get().Database.Collection(mutility.GetName(modelPtr))
	ctx, _ := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)

	res, err = collection.InsertOne(ctx, modelPtr)
	if err != nil {
		return nil, err
	}
	// modelPtr.ID = res.InsertedID
	val := reflect.ValueOf(modelPtr).Elem().FieldByName("ID")
	val.Set(reflect.ValueOf(res.InsertedID))
	return res, err
}

//InsertMany This will insert multiple Data
//TODO Find a way to pass pointer and attach its ID to the respective array elements
func InsertMany(models []interface{}) (res *mongo.InsertManyResult, err error) {
	if models == nil || len(models) == 0 {
		return nil, errors.New("The length of Model Array is 0")
	}
	collection := Get().Database.Collection(mutility.GetName(models))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	// iM := make([]interface{}, 0)
	// iM = append(iM, models)
	res, err = collection.InsertMany(ctx, models)
	if err != nil {
		return nil, err
	}
	return res, err
}
