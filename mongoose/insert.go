package mongoose

import (
	"context"
	"errors"
	"time"

	"github.com/rahul-sinha1908/go-mongoose/mutility"
	"go.mongodb.org/mongo-driver/mongo"
)

//InsertOne This will insert just one Data
func InsertOne(model interface{}) (res *mongo.InsertOneResult, err error) {
	collection := Get().Database.Collection(mutility.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	res, err = collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	return res, err
}

//InsertMany This will insert multiple Data
func InsertMany(models []interface{}) (res *mongo.InsertManyResult, err error) {
	if models == nil || len(models) == 0 {
		return nil, errors.New("The length of Model Array is 0")
	}
	collection := Get().Database.Collection(mutility.GetName(models))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	iM := make([]interface{}, 0)
	iM = append(iM, models)
	res, err = collection.InsertMany(ctx, iM)
	if err != nil {
		return nil, err
	}
	return res, err
}
