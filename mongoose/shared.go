package mongoose

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentBase struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}
