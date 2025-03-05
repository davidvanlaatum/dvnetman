package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CableConnectionType string

const (
	CableConnectionTypeDevice CableConnectionType = "device"
)

type CableConnection struct {
	Type *CableConnectionType `bson:"type,omitempty"`
	ID   *UUID                `bson:"id,omitempty"`
}

type Cable struct {
	Base        `bson:",inline"`
	Name        string             `bson:"name,omitempty"`
	Length      *float32           `bson:"length,omitempty"`
	Connections [2]CableConnection `bson:"connections"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "cables",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (c *Cable) GetCollectionName() string {
	return "cables"
}
