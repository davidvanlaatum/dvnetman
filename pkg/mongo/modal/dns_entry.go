package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DNSEntry struct {
	Base    `bson:",inline"`
	DNSZone *UUID  `bson:"dns_zone,omitempty"`
	Name    string `bson:"name,omitempty"`
	Type    string `bson:"type,omitempty"`
	Value   string `bson:"value,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "dns_entries",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *DNSEntry) GetCollectionName() string {
	return "dns_entries"
}
