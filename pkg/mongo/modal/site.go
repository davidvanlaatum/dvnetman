package modal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Site struct {
	Base        `bson:",inline"`
	Name        string  `bson:"name,omitempty"`
	Description *string `bson:"description,omitempty"`
	TimeZone    *string `bson:"timezone,omitempty"`
	Latitude    *string `bson:"latitude,omitempty"`
	Longitude   *string `bson:"longitude,omitempty"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "sites",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *Site) GetCollectionName() string {
	return "sites"
}

func (db *DBClient) GetSite(
	ctx context.Context, id *UUID, opts ...options.Lister[options.FindOneOptions],
) (d *Site, err error) {
	d = &Site{}
	err = findById(ctx, db, id, &d, opts...)
	return
}

func (db *DBClient) ListSites(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (sites []*Site, err error) {
	sites = []*Site{nil}
	err = listBy(ctx, db, filter, &sites, opts...)
	return
}

func (db *DBClient) SaveSite(ctx context.Context, site *Site) error {
	return save(ctx, db, &site)
}

func (db *DBClient) DeleteSite(ctx context.Context, site *Site) error {
	return deleteObj(ctx, db, &site)
}

func (db *DBClient) CountSites(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions],
) (int64, error) {
	t := &Site{}
	return count(ctx, db, &t, filter, opts...)
}
