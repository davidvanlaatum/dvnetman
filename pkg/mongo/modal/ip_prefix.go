package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IPPrefix struct {
	Base        `bson:",inline"`
	Prefix      string  `bson:"prefix,omitempty"`
	Name        *string `bson:"name,omitempty"`
	Description *string `bson:"description,omitempty"`
	VRFDomain   *UUID   `bson:"vrf_domain,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "ip_prefixes",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (i *IPPrefix) GetCollectionName() string {
	return "ip_prefixes"
}
