package mongoose

import "go.mongodb.org/mongo-driver/v2/bson"

type DocumentBase struct {
	ID bson.ObjectID `json:"_id" bson:"_id"`
}
