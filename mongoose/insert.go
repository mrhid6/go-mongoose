package mongoose

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/mrhid6/go-mongoose/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

// InsertOne This will insert just one Data
//
// Parameters:
//
//   - modelPtr - Should be a pointer to the model e.g. &NewUser
func InsertOne(modelPtr interface{}) (res *mongo.InsertOneResult, err error) {

	if !utils.IsPointer(modelPtr) {
		return nil, errors.New("insertone - model should be a Pointer")
	}

	mongo, err := Get()

	if err != nil {
		return nil, err
	}

	collection := mongo.Database.Collection(utils.GetName(modelPtr))
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
//
// Parameters:
//
//   - modelsPtr - Should be a pointer to the model array e.g. &[]NewUser
func InsertMany(modelsPtr []interface{}) (res *mongo.InsertManyResult, err error) {
	if len(modelsPtr) == 0 {
		return nil, errors.New("the length of Model Array is 0")
	}

	if !utils.IsPointer(modelsPtr) {
		return nil, errors.New("models should be a Pointer")
	}

	mongo, err := Get()

	if err != nil {
		return nil, err
	}

	collection := mongo.Database.Collection(utils.GetName(modelsPtr))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	// iM := make([]interface{}, 0)
	// iM = append(iM, models)
	res, err = collection.InsertMany(ctx, modelsPtr)
	if err != nil {
		return nil, err
	}
	return res, err
}
