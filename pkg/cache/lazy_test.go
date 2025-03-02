package cache

import (
	"context"
	"dvnetman/pkg/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLazy(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m, err := NewPool(ctx, "memory://?ttl=1s")
	r.NoError(err)
	defer func() {
		r.NoError(m.Shutdown(ctx))
	}()

	var called int
	for i := 0; i < 2; i++ {
		var v string
		v, err = Lazy(
			ctx, m, "key", func(ctx context.Context) (string, error) {
				called++
				return "value", nil
			},
		)
		r.NoError(err)
		r.Equal("value", v)
	}
	r.Equal(1, called)
}
