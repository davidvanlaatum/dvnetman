package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const versionField = "version"

type Base struct {
	ObjectId bson.ObjectID `bson:"_id,omitempty"`
	ID       *UUID         `bson:"id,minsize:36"`
	Created  time.Time     `bson:"created"`
	Updated  time.Time     `bson:"updated"`
	Version  int           `bson:"version"`
	Tags     []string      `bson:"tags,omitempty"`
	loaded   []byte        `bson:"-"`
}

func (b *Base) baseImpl() {
	// NOOP
}

func (b *Base) GetBase() *Base {
	return b
}

type baseInterface interface {
	baseImpl()
	GetBase() *Base
	GetCollectionName() string
}
