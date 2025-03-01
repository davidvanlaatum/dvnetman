package modal

import (
	"bytes"
	"context"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/mongo/adapt"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type CollectionInfo struct {
	Name    string
	Indexes []mongo.IndexModel
}

var collections []*CollectionInfo

func register(c *CollectionInfo) {
	collections = append(collections, c)
}

var OptimisticLockError = errors.New("optimistic lock error")

type DBClient struct {
	db      mongoadapt.MongoDatabase
	getNow  func() time.Time
	newUUID func() *UUID
}

type DBClientOption func(*DBClient)

func WithNowFunc(f func() time.Time) DBClientOption {
	return func(c *DBClient) {
		c.getNow = f
	}
}

func WithUUIDFunc(f func() *UUID) DBClientOption {
	return func(c *DBClient) {
		c.newUUID = f
	}
}

func NewDBClient(db mongoadapt.MongoDatabase, opt ...DBClientOption) *DBClient {
	c := &DBClient{
		db:     db,
		getNow: time.Now,
		newUUID: func() *UUID {
			if id, err := uuid.NewV7(); err != nil {
				panic(err)
			} else {
				return (*UUID)(&id)
			}
		},
	}
	for _, option := range opt {
		option(c)
	}
	return c
}

func (db *DBClient) Init(ctx context.Context) (err error) {
	var c mongoadapt.MongoCursor
	if c, err = db.db.ListCollections(ctx, bson.M{}); err != nil {
		return errors.Wrap(err, "failed to list collections")
	}
	defer utils.PropagateErrorContext(ctx, c.Close, &err, "failed to close collections result")
	var list []bson.M
	if err = c.All(ctx, &list); err != nil {
		return errors.Wrap(err, "failed to get collections")
	}
	exists := map[string]bool{}
	for _, info := range list {
		exists[info["name"].(string)] = true
	}
	for _, info := range collections {
		if err = db.initCollection(ctx, info, exists[info.Name]); err != nil {
			return err
		}
	}
	return nil
}

func (db *DBClient) initCollection(ctx context.Context, info *CollectionInfo, exists bool) (err error) {
	if !exists {
		logger.Ctx(ctx).Info().Key("collection", info.Name).Msg("Creating collection")
		if err = db.db.CreateCollection(ctx, info.Name); err != nil {
			return errors.Wrapf(err, "failed to create collection %s", info.Name)
		}
	}
	col := db.db.Collection(info.Name)
	var indexes []mongo.IndexSpecification
	if indexes, err = col.Indexes().ListSpecifications(ctx); err != nil {
		return errors.Wrapf(err, "failed to list indexes for collection %s", info.Name)
	}
	var indexMap []mongo.IndexModel
	if indexMap, err = utils.MapErr(indexes, convertIndexSpecificationToModel); err != nil {
		return err
	}
	for _, index := range info.Indexes {
		if _, found := utils.FindFirst(
			indexMap, func(i mongo.IndexModel) bool {
				return indexModelEqual(i, index)
			},
		); !found {
			var name string
			if name, err = col.Indexes().CreateOne(ctx, index); err != nil {
				return errors.Wrapf(err, "failed to create index for collection %s", info.Name)
			}
			logger.Ctx(ctx).Info().Key("collection", info.Name).Key("index", name).Msg("Created index")
		}
	}
	return nil
}

func convertIndexSpecificationToModel(i mongo.IndexSpecification) (mongo.IndexModel, error) {
	k := bson.D{}
	if err := bson.Unmarshal(i.KeysDocument, &k); err != nil {
		return mongo.IndexModel{}, err
	}
	model := mongo.IndexModel{
		Keys: bson.D(
			utils.MapTo(
				k, func(kv bson.E) bson.E {
					return bson.E{Key: kv.Key, Value: kv.Value}
				},
			),
		),
		Options: options.Index(),
	}
	if i.Unique != nil {
		model.Options = model.Options.SetUnique(*i.Unique)
	}
	return model, nil
}

func indexOptionsBuilderToIndexOptions(i *options.IndexOptionsBuilder) options.IndexOptions {
	var opts options.IndexOptions
	if i != nil {
		for _, f := range i.List() {
			if err := f(&opts); err != nil {
				panic(err)
			}
		}
	}
	return opts
}

func indexModelEqual(a, b mongo.IndexModel) bool {
	if len(a.Keys.(bson.D)) != len(b.Keys.(bson.D)) {
		return false
	}
	for i, v := range a.Keys.(bson.D) {
		if v.Key != b.Keys.(bson.D)[i].Key || v.Value != b.Keys.(bson.D)[i].Value {
			return false
		}
	}

	aopts := indexOptionsBuilderToIndexOptions(a.Options)
	bopts := indexOptionsBuilderToIndexOptions(b.Options)

	return utils.ComparePointers(aopts.Unique, bopts.Unique)
}

func findById[T baseInterface](
	ctx context.Context, db *DBClient, id *UUID, result *T, opts ...options.Lister[options.FindOneOptions],
) error {
	err := collection(db, *result).FindOne(ctx, bson.M{"id": id}, opts...).Decode(&result)
	if err != nil {
		result = nil
		return errors.Wrapf(err, "failed to find document with id %s", id)
	}
	(*result).GetBase().loaded, err = bson.Marshal(*result)
	if err != nil {
		return errors.Wrap(err, "failed to marshal document")
	}
	return nil
}

func listBy[T baseInterface](
	ctx context.Context, db *DBClient, filter interface{}, result *[]T, opts ...options.Lister[options.FindOptions],
) error {
	c, err := collection(db, (*result)[0]).Find(ctx, filter, opts...)
	if err != nil {
		*result = nil
		return errors.Wrap(err, "failed to find document")
	}
	*result = (*result)[:0]
	if err = c.All(ctx, result); err != nil {
		return errors.Wrap(err, "failed to get documents")
	}
	for _, r := range *result {
		r.GetBase().loaded, err = bson.Marshal(r)
		if err != nil {
			return errors.Wrap(err, "failed to marshal document")
		}
	}
	return nil
}

func findOne[T baseInterface](
	ctx context.Context, db *DBClient, filter interface{}, result *T, opts ...options.Lister[options.FindOneOptions],
) error {
	err := collection(db, *result).FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		result = nil
		return errors.Wrap(err, "failed to find document")
	}
	(*result).GetBase().loaded, err = bson.Marshal(*result)
	if err != nil {
		return errors.Wrap(err, "failed to marshal document")
	}
	return nil
}

func collection(db *DBClient, obj baseInterface) mongoadapt.MongoCollection {
	return db.db.Collection(obj.GetCollectionName())
}

func save[T baseInterface](ctx context.Context, db *DBClient, document *T) (err error) {
	now := db.getNow().UTC().Truncate(time.Millisecond)
	base := (*document).GetBase()

	if base.ObjectId.IsZero() {
		base.Created = now
		base.Updated = now
		base.ID = db.newUUID()
		base.Version = 1
		var res *mongo.InsertOneResult
		if res, err = collection(db, *document).InsertOne(ctx, document); err != nil {
			return errors.Wrap(err, "failed to insert document")
		}
		base.ObjectId = res.InsertedID.(bson.ObjectID)
		base.loaded, err = bson.Marshal(*document)
	} else if hasChanged(document) {
		base.Updated = now
		var res *mongo.UpdateResult
		oldVersion := base.Version
		base.Version++
		if res, err = collection(db, *document).UpdateOne(
			ctx, bson.M{
				"_id":        base.ObjectId,
				versionField: oldVersion,
			}, bson.M{"$set": document},
		); err != nil {
			base.Version = oldVersion
			return errors.Wrap(err, "failed to update document")
		} else if res.MatchedCount == 0 {
			base.Version = oldVersion
			return errors.WithStack(OptimisticLockError)
		}
		base.loaded, err = bson.Marshal(*document)
	}
	return
}

func hasChanged[T baseInterface](document *T) bool {
	loaded := (*document).GetBase().loaded
	if loaded == nil {
		return true
	}
	current, err := bson.Marshal(*document)
	if err != nil {
		return true
	}
	return !bytes.Equal(current, loaded)
}

func deleteObj[T baseInterface](ctx context.Context, db *DBClient, document *T) error {
	base := (*document).GetBase()
	del, err := collection(db, *document).DeleteOne(
		ctx, bson.M{
			"_id":        base.ObjectId,
			versionField: base.Version,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to delete document")
	}
	if del.DeletedCount == 0 {
		return errors.WithStack(OptimisticLockError)
	}
	return nil
}

func count[T baseInterface](
	ctx context.Context, db *DBClient, document *T, filter interface{}, opts ...options.Lister[options.CountOptions],
) (int64, error) {
	c, err := collection(db, *document).CountDocuments(ctx, filter, opts...)
	return c, errors.Wrap(err, "failed to count documents")
}
