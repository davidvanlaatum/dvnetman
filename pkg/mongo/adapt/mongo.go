package mongoadapt

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoCursor see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Cursor
type MongoCursor interface {
	// Next see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Cursor.Next
	Next(ctx context.Context) bool
	// Decode see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Cursor.Decode
	Decode(v interface{}) error
	// Close see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Cursor.Close
	Close(ctx context.Context) error
	// All see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Cursor.All
	All(ctx context.Context, result interface{}) error
}

// MongoDatabase see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Database
type MongoDatabase interface {
	ListCollections(
		ctx context.Context, filter interface{}, opts ...options.Lister[options.ListCollectionsOptions],
	) (MongoCursor, error)
	CreateCollection(ctx context.Context, name string, opts ...options.Lister[options.CreateCollectionOptions]) error
	Collection(name string, opts ...options.Lister[options.CollectionOptions]) MongoCollection
}

func AdapterMongoDatabase(db *mongo.Database) MongoDatabase {
	return &mongoDatabase{Database: db}
}

type MongoCollection interface {
	// InsertOne see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.InsertOne
	InsertOne(
		ctx context.Context, document interface{}, opts ...options.Lister[options.InsertOneOptions],
	) (*mongo.InsertOneResult, error)
	// Indexes see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.Indexes
	Indexes() MongoIndexView
	// UpdateOne see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.UpdateOne
	UpdateOne(
		ctx context.Context, filter interface{}, update interface{}, opts ...options.Lister[options.UpdateOneOptions],
	) (*mongo.UpdateResult, error)
	// ReplaceOne see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.ReplaceOne
	ReplaceOne(
		ctx context.Context, filter interface{}, update interface{}, opts ...options.Lister[options.ReplaceOptions],
	) (*mongo.UpdateResult, error)
	// DeleteOne see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.DeleteOne
	DeleteOne(
		ctx context.Context, filter interface{}, opts ...options.Lister[options.DeleteOneOptions],
	) (*mongo.DeleteResult, error)
	// FindOne see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.FindOne
	FindOne(ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	// Find see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.Find
	Find(ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions]) (MongoCursor, error)
	// CountDocuments see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.CountDocuments
	CountDocuments(ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions]) (int64, error)
}

type MongoIndexView interface {
	ListSpecifications(
		ctx context.Context, opts ...options.Lister[options.ListIndexesOptions],
	) ([]mongo.IndexSpecification, error)
	CreateOne(
		ctx context.Context, model mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions],
	) (string, error)
}

type mongoDatabase struct {
	*mongo.Database
}

func (m *mongoDatabase) ListCollections(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.ListCollectionsOptions],
) (MongoCursor, error) {
	return m.Database.ListCollections(ctx, filter, opts...)
}

func (m *mongoDatabase) Collection(name string, opts ...options.Lister[options.CollectionOptions]) MongoCollection {
	return &mongoCollection{Collection: m.Database.Collection(name, opts...)}
}

var _ MongoDatabase = (*mongoDatabase)(nil)

type mongoCollection struct {
	*mongo.Collection
}

func (m *mongoCollection) Indexes() MongoIndexView {
	return m.Collection.Indexes()
}

func (m *mongoCollection) Find(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (MongoCursor, error) {
	return m.Collection.Find(ctx, filter, opts...)
}
