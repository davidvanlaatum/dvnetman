package modal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	Base             `bson:",inline"`
	FirstName        *string `bson:"first_name"`
	LastName         *string `bson:"last_name"`
	FullName         *string `bson:"full_name"`
	Email            *string `bson:"email"`
	Password         *string `bson:"password"`
	ExternalProvider *string `bson:"external_provider"`
	ExternalID       *string `bson:"external_id"`
	Username         *string `bson:"username"`
}

func (u *User) GetCollectionName() string {
	return "users"
}

func init() {
	register(
		&CollectionInfo{
			Name: "users",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"id", int32(1)}},
					Options: options.Index().SetUnique(true),
				},
			},
		},
	)
}

func (db *DBClient) GetUser(ctx context.Context, m *UUID, opts ...options.Lister[options.FindOneOptions]) (
	*User, error,
) {
	u := &User{}
	err := findById(ctx, db, m, &u, opts...)
	return u, err
}

func (db *DBClient) SaveUser(ctx context.Context, mod *User) error {
	return save(ctx, db, &mod)
}

func (db *DBClient) DeleteUser(ctx context.Context, d *User) error {
	return deleteObj(ctx, db, &d)
}

func (db *DBClient) ListUsers(ctx context.Context, search interface{}, opts ...options.Lister[options.FindOptions]) (
	users []*User, err error,
) {
	users = []*User{nil}
	err = listBy(ctx, db, search, &users, opts...)
	return
}

func (db *DBClient) CountUsers(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions],
) (int64, error) {
	t := &User{}
	return count(ctx, db, &t, filter, opts...)
}
