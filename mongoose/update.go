package mongoose

import (
	"context"
	"time"

	"github.com/rahul-sinha1908/go-mongoose/interfaces"
	"github.com/rahul-sinha1908/go-mongoose/mutility"
	"go.mongodb.org/mongo-driver/bson"
)

//UpdateByID Updates by ID
func UpdateByID(model interfaces.ModelInterface) error {
	collection := Get().Database.Collection(mutility.GetName(model))
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	_, err := collection.UpdateOne(ctx, bson.M{
		"_id": mutility.GetID(model),
	}, model)

	if err != nil {
		return err
	}
	return nil
}
