package modal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DeviceTypePort struct {
	ID             string       `bson:"id,omitempty"`
	Name           string       `bson:"name,omitempty"`
	Label          string       `bson:"label,omitempty"`
	PortType       PortType     `bson:"port_type,omitempty"`
	PortLocation   PortLocation `bson:"port_location,omitempty"`
	Type           string       `bson:"type,omitempty"`
	Enabled        *bool        `bson:"enabled,omitempty"`
	ManagementOnly *bool        `bson:"management_only,omitempty"`
	Description    string       `bson:"description,omitempty"`
	Bridge         string       `bson:"bridge,omitempty"`
	PoEMode        PoEMode      `bson:"poe_mode,omitempty"`
	PoEType        PoEType      `bson:"poe_type,omitempty"`
}

type DeviceType struct {
	Base         `bson:",inline"`
	Manufacturer *UUID  `bson:"manufacturer,omitempty"`
	Model        string `bson:"model,omitempty"`
	Description  string `bson:"description,omitempty"`
	Height       int    `bson:"height,omitempty"`
	FullDepth    bool   `bson:"full_depth,omitempty"`
	PartNumber   string `bson:"part_number,omitempty"`
	Comments     string `bson:"comments,omitempty"`

	Ports []DeviceTypePort `bson:"ports,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "device_types",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *DeviceType) GetCollectionName() string {
	return "device_types"
}

func (db *DBClient) GetDeviceType(
	ctx context.Context, id *UUID, opts ...options.Lister[options.FindOneOptions],
) (d *DeviceType, err error) {
	d = &DeviceType{}
	err = findById(ctx, db, id, &d, opts...)
	return
}

func (db *DBClient) ListDeviceTypes(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (devices []*DeviceType, err error) {
	devices = []*DeviceType{nil}
	err = listBy(ctx, db, filter, &devices, opts...)
	return
}

func (db *DBClient) SaveDeviceType(ctx context.Context, device *DeviceType) error {
	return save(ctx, db, &device)
}

func (db *DBClient) DeleteDeviceType(ctx context.Context, device *DeviceType) error {
	return deleteObj(ctx, db, &device)
}
