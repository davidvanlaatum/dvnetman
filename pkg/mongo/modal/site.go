package modal

import "go.mongodb.org/mongo-driver/v2/bson"

type Site struct {
	ObjectId    bson.ObjectID `bson:"_id,omitempty"`
	ID          string        `bson:"_id,omitempty"`
	Name        string        `bson:"name,omitempty"`
	Description string        `bson:"description,omitempty"`
	TimeZone    string        `bson:"timezone,omitempty"`
	Latitude    string        `bson:"latitude,omitempty"`
	Longitude   string        `bson:"longitude,omitempty"`
}
