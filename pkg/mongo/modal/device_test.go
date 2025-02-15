package modal

import (
	"dvnetman/pkg/mongo/adapt"
	"dvnetman/pkg/testutils"
	"dvnetman/pkg/utils"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"testing"
	"time"
)

func TestDevice(t *testing.T) {
	m := testutils.StartMongo(t, testutils.GetTestContext(t))
	t.Run(
		"SaveAndGetDevice", func(t *testing.T) {
			r := require.New(t)
			ctx := testutils.GetTestContext(t)
			client := testutils.GetMongoClient(t, ctx, m)
			db := testutils.GetTestCaseDatabase(t, client)
			c := NewDBClient(mongoadapt.AdapterMongoDatabase(db))
			r.NoError(c.Init(ctx))
			r.NoError(c.Init(ctx))
			dev := &Device{
				Name: utils.ToPtr("Test"),
			}
			err := c.SaveDevice(ctx, dev)
			r.NoError(err)
			r.NotEmpty(dev.ID)
			r.NotEmpty(dev.ObjectId)
			r.NotEmpty(dev.Created)
			r.NotEmpty(dev.Updated)
			r.Equal(dev.Created, dev.Updated)
			dev2, err := c.GetDevice(ctx, dev.ID)
			r.NoError(err)
			r.Equal(dev.ID, dev2.ID)
			r.EqualExportedValues(dev, dev2)
			t.Logf("Device: %+v", dev2)

			r.NoError(c.SaveDevice(ctx, dev2))
			dev2.Name = utils.ToPtr("Test2")
			time.Sleep(time.Millisecond)
			r.NoError(c.SaveDevice(ctx, dev2))
			r.NotEqual(dev.Version, dev2.Version)
			r.NotEqual(dev.Created, dev2.Updated)

			dev3 := map[string]interface{}{}
			r.NoError(db.Collection(dev.GetCollectionName()).FindOne(ctx, bson.M{"id": dev.ID}).Decode(dev3))
			t.Logf("Device: %+v", dev3)

			dev.Name = utils.ToPtr("Test3")
			r.Error(c.SaveDevice(ctx, dev))
			r.Error(c.DeleteDevice(ctx, dev))

			r.NoError(c.DeleteDevice(ctx, dev2))
		},
	)

	t.Run(
		"ListDevices", func(t *testing.T) {
			r := require.New(t)
			ctx := testutils.GetTestContext(t)
			client := testutils.GetMongoClient(t, ctx, m)
			c := NewDBClient(mongoadapt.AdapterMongoDatabase(testutils.GetTestCaseDatabase(t, client)))
			r.NoError(c.Init(ctx))

			var devices []*Device
			for i := 0; i < 10; i++ {
				dev := &Device{
					Name: utils.ToPtr(fmt.Sprintf("Test%d", i)),
				}
				r.NoError(c.SaveDevice(ctx, dev))
				devices = append(devices, dev)
			}

			devices2, err := c.ListDevices(ctx, bson.M{})
			r.NoError(err)
			r.Len(devices2, 10)
			devices2, err = c.ListDevices(
				ctx, bson.M{"name": bson.M{"$in": []string{"Test2", "Test1", "Test3", "Test4"}}}, options.Find().
					SetSort(bson.D{{"name", -1}}).
					SetLimit(1).
					SetSkip(2),
			)
			r.NoError(err)
			r.Len(devices2, 1)
			r.Equal(devices[2], devices2[0])
		},
	)
}
