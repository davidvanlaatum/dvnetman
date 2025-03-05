package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IPAddress struct {
	Base        `bson:",inline"`
	IP          string   `bson:"ip,omitempty"`
	Device      *UUID    `bson:"device,omitempty"`
	Description *string  `bson:"description,omitempty"`
	DNSNames    []string `bson:"dns_name,omitempty"`
	VRFDomain   *UUID    `bson:"vrf_domain,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "ip_addresses",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (i *IPAddress) GetCollectionName() string {
	return "ip_addresses"
}
