package modal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type PoEMode string

const (
	PoEModePD  PoEMode = "PD"  // Powered device
	PoEModePSE PoEMode = "PSE" // Power-supplying equipment
)

type PoEType string

const (
	PoEType8023afT1     PoEType = "802.3af (Type 1)"
	PoEType8023atT2     PoEType = "802.3at (Type 2)"
	PoEType8023btT3     PoEType = "802.3bt (Type 3)"
	PoEType8023btT4     PoEType = "802.3bt (Type 4)"
	PoETypePassive24V2P PoEType = "Passive 24V 2-pair"
	PoETypePassive24V4P PoEType = "Passive 24V 4-pair"
	PoETypePassive48V2P PoEType = "Passive 48V 2-pair"
	PoETypePassive48V4P PoEType = "Passive 48V 4-pair"
)

type PortType string

const (
	PortTypeInterface     PortType = "interface"
	PortTypePowerIn       PortType = "power_in"
	PortTypePowerOut      PortType = "power_out"
	PortTypeConsole       PortType = "console"
	PortTypeConsoleServer PortType = "console_server"
)

type PortLocation string

const (
	PortLocationFront   PortLocation = "front"
	PortLocationRear    PortLocation = "rear"
	PortLocationVirtual PortLocation = "virtual"
)

type DevicePortType string

const (
	DevicePortTypeVirtual DevicePortType = "virtual"
	DevicePortTypeBridge  DevicePortType = "bridge"
	DevicePortTypeLAG     DevicePortType = "LAG"
)

type VLANMode string

const (
	VLANModeAccess    VLANMode = "access"
	VLANModeTagged    VLANMode = "tagged"
	VLANModeTaggedAll VLANMode = "tagged-all"
	VLANModeQinQ      VLANMode = "qinq" // TODO
)

type DevicePortVLANTaggedRange struct {
	Start int  `bson:"start,omitempty"`
	End   *int `bson:"end,omitempty"`
}

type DevicePort struct {
	Id             string                      `bson:"id,minsize:36"`
	Name           string                      `bson:"name"`
	Label          string                      `bson:"label,omitempty"`
	Description    string                      `bson:"description,omitempty"`
	Tags           []string                    `bson:"tags,omitempty"`
	PortType       PortType                    `bson:"port_type,omitempty"`
	PortLocation   PortLocation                `bson:"port_location,omitempty"`
	Type           DevicePortType              `bson:"type,omitempty"`
	Enabled        *bool                       `bson:"enabled,omitempty"`
	ManagementOnly *bool                       `bson:"management_only,omitempty"`
	Bridge         string                      `bson:"bridge,omitempty"`
	Parent         string                      `bson:"parent,omitempty"`
	LAG            string                      `bson:"lag,omitempty"`
	PoEMode        PoEMode                     `bson:"poe_mode,omitempty"`
	PoEType        PoEType                     `bson:"poe_type,omitempty"`
	VLANMode       VLANMode                    `bson:"vlan_mode,omitempty"`
	VLANGroup      string                      `bson:"vlan_group,omitempty"`
	VLANUntagged   int                         `bson:"vlan_untagged,omitempty"`
	VLANTagged     []DevicePortVLANTaggedRange `bson:"vlan_tagged,omitempty"`
}

type DeviceStatus string

const (
	DeviceStatusOffline         DeviceStatus = "offline"
	DeviceStatusActive          DeviceStatus = "active"
	DeviceStatusPlanned         DeviceStatus = "planned"
	DeviceStatusStaged          DeviceStatus = "staged"
	DeviceStatusFailed          DeviceStatus = "failed"
	DeviceStatusInventory       DeviceStatus = "inventory"
	DeviceStatusDecommissioning DeviceStatus = "decommissioning"
)

type RackFace string

const (
	RackFaceFront RackFace = "front"
	RackFaceBack  RackFace = "back"
)

type Slot struct {
	Id   string `bson:"id,minsize:36"`
	Name string `bson:"name"`
	Type string `bson:"type,omitempty"`
}

type Device struct {
	Base         `bson:",inline"`
	Name         *string       `bson:"name,omitempty"`
	Description  *string       `bson:"description,omitempty"`
	Status       *DeviceStatus `bson:"status,omitempty"`
	Site         *UUID         `bson:"site,omitempty"`
	Location     *UUID         `bson:"location,omitempty"`
	Rack         *UUID         `bson:"rack,omitempty"`
	RackFace     *RackFace     `bson:"rack_face,omitempty"`
	RackPosition *int          `bson:"position,omitempty"`
	DeviceType   *UUID         `bson:"device_type,omitempty"`
	Serial       *string       `bson:"serial,omitempty"`
	AssetTag     *string       `bson:"asset_tag,omitempty"`
	InstallDate  *time.Time    `bson:"install_date,omitempty"`

	Ports []DevicePort `bson:"ports,omitempty"`
	Slots []Slot       `bson:"slots,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "devices",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *Device) GetCollectionName() string {
	return "devices"
}

func (db *DBClient) GetDevice(
	ctx context.Context, id *UUID, opts ...options.Lister[options.FindOneOptions],
) (d *Device, err error) {
	d = &Device{}
	err = findById(ctx, db, id, &d, opts...)
	return
}

func (db *DBClient) ListDevices(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (devices []*Device, err error) {
	devices = []*Device{nil}
	err = listBy(ctx, db, filter, &devices, opts...)
	return
}

func (db *DBClient) SaveDevice(ctx context.Context, device *Device) error {
	return save(ctx, db, &device)
}

func (db *DBClient) DeleteDevice(ctx context.Context, device *Device) error {
	return deleteObj(ctx, db, &device)
}

func (db *DBClient) CountDevices(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions],
) (int64, error) {
	t := &Device{}
	return count(ctx, db, &t, filter, opts...)
}
