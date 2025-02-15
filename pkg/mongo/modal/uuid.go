package modal

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UUID uuid.UUID

func (u *UUID) MarshalBSONValue() (byte, []byte, error) {
	t, b, err := bson.MarshalValue((*uuid.UUID)(u).String())
	return byte(t), b, err
}

func (u *UUID) UnmarshalBSONValue(typ byte, data []byte) error {
	var x string
	if err := bson.UnmarshalValue(bson.Type(typ), data, &x); err != nil {
		return err
	}
	uu, err := uuid.Parse(x)
	if err != nil {
		return err
	}
	*u = UUID(uu)
	return nil
}

var _ bson.ValueMarshaler = (*UUID)(nil)
var _ bson.ValueUnmarshaler = (*UUID)(nil)
