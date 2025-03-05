package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FHRPGroupProtocol string

const (
	FHRPGroupProtocolHSRP FHRPGroupProtocol = "HSRP"
	FHRPGroupProtocolVRRP FHRPGroupProtocol = "VRRP"
	FHRPGroupProtocolGLBP FHRPGroupProtocol = "GLBP"
	FHRPGroupProtocolCARP FHRPGroupProtocol = "CARP"
)

type FHRPGroup struct {
	Base        `bson:",inline"`
	Name        string            `bson:"name,omitempty"`
	Description string            `bson:"description,omitempty"`
	GroupID     int               `bson:"groupNumber,omitempty"`
	Protocol    FHRPGroupProtocol `bson:"protocol,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "fhrp_groups",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (f *FHRPGroup) GetCollectionName() string {
	return "fhrp_groups"
}
