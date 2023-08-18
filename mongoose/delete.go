package mongoose

import (
	"context"
	"time"

	"github.com/mrhid6/go-mongoose/mutility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//DeleteOne Deletes one Object
func DeleteOne(filter bson.M, tempCollection interface{}) (*mongo.DeleteResult, error) {
	// fmt.Println("Collection Name : ", b.GetName())
	collection := Get().Database.Collection(mutility.GetName(tempCollection))
	ctx, _ := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)

	return collection.DeleteOne(ctx, filter)
}

//DeleteMany Deletes Many Objects
func DeleteMany(filter bson.M, tempCollection interface{}) (*mongo.DeleteResult, error) {
	// fmt.Println("Collection Name : ", b.GetName())
	collection := Get().Database.Collection(mutility.GetName(tempCollection))
	ctx, _ := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)

	return collection.DeleteMany(ctx, filter)
}
