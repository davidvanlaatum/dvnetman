package modal

import "go.mongodb.org/mongo-driver/v2/bson"

type Location struct {
	ObjectId    bson.ObjectID `bson:"_id,omitempty"`
	ID          string        `bson:"id,omitempty"`
	Parent      string        `bson:"parent,omitempty"`
	Name        string        `bson:"name,omitempty"`
	Description string        `bson:"description,omitempty"`
}
