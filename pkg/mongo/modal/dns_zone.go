package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DNSZone struct {
	Base        `bson:",inline"`
	Name        string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "dns_zones",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *DNSZone) GetCollectionName() string {
	return "dns_zones"
}
