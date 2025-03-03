package modal

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
