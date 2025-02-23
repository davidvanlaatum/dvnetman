package modal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Manufacturer struct {
	Base `bson:",inline"`
	Name string `bson:"name"`
}

func init() {
	register(
		&CollectionInfo{
			Name: "manufacturers",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys:    bson.D{{"name", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (d *Manufacturer) GetCollectionName() string {
	return "manufacturers"
}

func (db *DBClient) GetManufacturer(
	ctx context.Context, id *UUID, opts ...options.Lister[options.FindOneOptions],
) (d *Manufacturer, err error) {
	d = &Manufacturer{}
	err = findById(ctx, db, id, &d, opts...)
	return
}

func (db *DBClient) ListManufacturers(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (manufacturers []*Manufacturer, err error) {
	manufacturers = []*Manufacturer{nil}
	err = listBy(ctx, db, filter, &manufacturers, opts...)
	return
}

func (db *DBClient) SaveManufacturer(ctx context.Context, manufacturer *Manufacturer) error {
	return save(ctx, db, &manufacturer)
}

func (db *DBClient) DeleteManufacturer(ctx context.Context, manufacturer *Manufacturer) error {
	return deleteObj(ctx, db, &manufacturer)
}

func (db *DBClient) CountManufacturers(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions],
) (int64, error) {
	t := &Manufacturer{}
	return count(ctx, db, &t, filter, opts...)
}
