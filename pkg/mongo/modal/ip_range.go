package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IPRange struct {
	Base      `bson:",inline"`
	Start     string `bson:"start"`
	End       string `bson:"end"`
	VRFDomain *UUID  `bson:"vrf_domain,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "ip_ranges",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (i *IPRange) GetCollectionName() string {
	return "ip_ranges"
}
