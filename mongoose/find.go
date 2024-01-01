package mongoose

import (
	"context"
	"time"

	"github.com/mrhid6/go-mongoose/mutility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne Searches one object and returns its value
func FindOne(filter bson.M, b interface{}) (err error) {
	// fmt.Println("Collection Name : ", mutility.GetName(b))
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(b))
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
	// fmt.Println("Collection Name : ", b.GetName())
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(b))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res := collection.FindOne(ctx, bson.M{
		"_id": userID,
	})
	if res.Err() != nil {
		return res.Err()
	}
	err = res.Decode(b)
	if err != nil {
		return err
	}

	return nil
}

// FindByObjectID Searches by Object ID
func FindByObjectID(objectID primitive.ObjectID, bPtr interface{}) (err error) {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(bPtr))
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	res := collection.FindOne(ctx, bson.M{
		"_id": objectID,
	})
	if res.Err() != nil {
		return res.Err()
	}
	err = res.Decode(bPtr)
	if err != nil {
		return err
	}

	return nil
}

func findByObjectID(objectID primitive.ObjectID, collectionName string) (interface{}, error) {
	mongo, err := Get()

	if err != nil {
		return nil, err
	}

	collection := mongo.Database.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)

	res := collection.FindOne(ctx, bson.M{
		"_id": objectID,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var b interface{}
	err = res.Decode(&b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// FindAll Get All Docs
func FindAll(filter bson.M, modelsOutArrayPtr interface{}) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(modelsOutArrayPtr))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	err = cur.All(ctx, modelsOutArrayPtr)
	if err != nil {
		return err
	}
	return nil
}

// FindAllWithPagination Get All Docs with Pagination
func FindAllWithPagination(filter bson.M, start int64, count int64, modelsOutArrayPtr interface{}) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(modelsOutArrayPtr))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	cur, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &start,
		Limit: &count,
	})
	if err != nil {
		return err
	}
	err = cur.All(ctx, modelsOutArrayPtr)
	if err != nil {
		return err
	}
	return nil
}

// FindAllWithOptions Find all with options
func FindAllWithOptions(filter bson.M, option options.FindOptions, modelsOutArrayPtr interface{}) error {
	mongo, err := Get()

	if err != nil {
		return err
	}

	collection := mongo.Database.Collection(mutility.GetName(modelsOutArrayPtr))
	ctx, _ := context.WithTimeout(context.Background(), LongWaitTime*time.Second)

	cur, err := collection.Find(ctx, filter, &option)
	if err != nil {
		return err
	}
	err = cur.All(ctx, modelsOutArrayPtr)
	if err != nil {
		return err
	}
	return nil
}
