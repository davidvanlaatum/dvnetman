package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type VRFDomain struct {
	Base        `bson:",inline"`
	Name        string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "vrf_domains",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (v *VRFDomain) GetCollectionName() string {
	return "vrf_domains"
}
