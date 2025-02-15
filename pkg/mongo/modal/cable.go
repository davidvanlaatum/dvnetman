package modal

import "go.mongodb.org/mongo-driver/v2/bson"

type CableConnectionType string

const (
	CableConnectionTypeDevice CableConnectionType = "device"
)

type CableConnection struct {
	Type CableConnectionType `bson:"type"`
	ID   bson.ObjectID       `bson:"id"`
}

type Cable struct {
	ID          bson.ObjectID      `bson:"_id"`
	CableID     string             `bson:"cable-id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Length      float32            `bson:"length,omitempty"`
	Connections [2]CableConnection `bson:"connections"`
}
