package modal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Location struct {
	Base        `bson:",inline"`
	Site        *UUID   `bson:"site,omitempty"`
	Parent      *UUID   `bson:"parent,omitempty"`
	Name        string  `bson:"name,omitempty"`
	Description *string `bson:"description,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "locations",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *Location) GetCollectionName() string {
	return "locations"
}

func (db *DBClient) GetLocation(
	ctx context.Context, id *UUID, opts ...options.Lister[options.FindOneOptions],
) (d *Location, err error) {
	d = &Location{}
	err = findById(ctx, db, id, &d, opts...)
	return
}

func (db *DBClient) ListLocations(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (locations []*Location, err error) {
	locations = []*Location{nil}
	err = listBy(ctx, db, filter, &locations, opts...)
	return
}

func (db *DBClient) SaveLocation(ctx context.Context, location *Location) error {
	return save(ctx, db, &location)
}

func (db *DBClient) DeleteLocation(ctx context.Context, location *Location) error {
	return deleteObj(ctx, db, &location)
}

func (db *DBClient) CountLocations(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions],
) (int64, error) {
	t := &Location{}
	return count(ctx, db, &t, filter, opts...)
}
