package modal

import "go.mongodb.org/mongo-driver/v2/bson"

type Rack struct {
	ObjectId bson.ObjectID `bson:"_id,omitempty"`
	ID       string        `bson:"id,omitempty"`
	Name     string        `bson:"name,omitempty"`
	Location string        `bson:"location,omitempty"`
	Height   int           `bson:"height,omitempty"`
	Width    int           `bson:"width,omitempty"`
}
