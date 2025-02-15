package modal

import "go.mongodb.org/mongo-driver/v2/bson"

type IPRange struct {
	ObjectID bson.ObjectID `bson:"_id,omitempty"`
	ID       string        `bson:"id,omitempty"`
	Start    string        `bson:"start"`
	End      string        `bson:"end"`
}
