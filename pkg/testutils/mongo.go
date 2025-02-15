package testutils

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"strings"
	"testing"
	"time"
)

func StartMongo(t testing.TB, ctx context.Context) *mongodb.MongoDBContainer {
	m, err := mongodb.Run(ctx, "mongo:latest")
	r := require.New(t)
	r.NoError(err)
	t.Cleanup(func() {
		testcontainers.CleanupContainer(t, m)
	})
	return m
}

func GetMongoClient(t testing.TB, ctx context.Context, m *mongodb.MongoDBContainer) *mongo.Client {
	r := require.New(t)
	endpoint, err := m.ConnectionString(ctx)
	r.NoError(err)
	client, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	r.NoError(err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r.NoError(client.Disconnect(ctx))
	})
	return client
}

func GetTestCaseDatabase(t testing.TB, client *mongo.Client) *mongo.Database {
	name := strings.ReplaceAll(t.Name(), "/", "_")
	t.Logf("using database %s", name)
	return client.Database(name)
}
