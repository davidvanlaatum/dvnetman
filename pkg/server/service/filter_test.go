package service

import (
	"dvnetman/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

func TestFilter(t *testing.T) {
	r := require.New(t)
	f := filter{}
	f.equalsStr("field", utils.ToPtr("value"))
	r.Len(f, 1)
	r.Equal(bson.E{Key: "field", Value: "value"}, f[0])
	f.equalsStr("field2", utils.ToPtr("value2"))
	f.regex("field3", utils.ToPtr(".*"), "i")
	b, err := bson.MarshalExtJSON(f, false, false)
	r.NoError(err)
	r.JSONEq(
		`{"field":"value","field2":"value2", "field3":{"$regularExpression":{"options":"i", "pattern":".*"}}}`, string(b),
	)
}
