package cache

import (
	"dvnetman/pkg/testutils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMemoryLookup(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m, err := NewPool(ctx, "memory://?ttl=1s")
	r.NoError(err)
	defer func() {
		r.NoError(m.Shutdown(ctx))
	}()
	r.NoError(m.Set(ctx, "key", "value"))
	var rt string
	var ok bool
	ok, err = m.Get(ctx, "key", &rt)
	r.NoError(err)
	r.True(ok)
	r.Equal("value", rt)
}

func TestMemoryNotFound(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m, err := NewPool(ctx, "memory://?ttl=1s")
	r.NoError(err)
	defer func() {
		r.NoError(m.Shutdown(ctx))
	}()

	var rt string
	var ok bool
	ok, err = m.Get(ctx, "key2", &rt)
	r.NoError(err)
	r.False(ok)
	r.Empty(rt)
}

func TestMemoryNew(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m := &MemoryCache{}
	m2, err := m.New(ctx)
	r.NoError(err)
	r.Same(m, m2)
}

func TestMemoryDelete(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m, err := NewPool(ctx, "memory://?ttl=1s")
	r.NoError(err)
	defer func() {
		r.NoError(m.Shutdown(ctx))
	}()
	r.NoError(m.Set(ctx, "key", "value"))
	r.NoError(m.Delete(ctx, "key"))
	var rt string
	var ok bool
	ok, err = m.Get(ctx, "key", &rt)
	r.NoError(err)
	r.False(ok)
}

func TestMemoryLimit(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m, err := NewPool(ctx, "memory://?ttl=2s&limit=15b")
	r.NoError(err)
	defer func() {
		r.NoError(m.Shutdown(ctx))
	}()

	r.NoError(m.Set(ctx, "key", "value"))
	r.NoError(m.Set(ctx, "key2", "value2"))
	r.Eventually(
		func() bool {
			var tmp string
			ok, _ := m.Get(ctx, "key", &tmp)
			return !ok
		}, time.Second, 10*time.Millisecond,
	)
	var rt string
	var ok bool
	ok, err = m.Get(ctx, "key", &rt)
	r.NoError(err)
	r.False(ok)
	ok, err = m.Get(ctx, "key2", &rt)
	r.NoError(err)
	r.True(ok)
	r.Equal("value2", rt)
}

func TestMemoryCleanup(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	m, err := NewPool(ctx, "memory://?ttl=10ms")
	r.NoError(err)
	defer func() {
		r.NoError(m.Shutdown(ctx))
	}()

	r.NoError(m.Set(ctx, "key", "value"))
	r.Eventually(
		func() bool {
			var tmp string
			ok, _ := m.Get(ctx, "key", &tmp)
			return !ok
		}, time.Second, 10*time.Millisecond,
	)
	var rt string
	var ok bool
	ok, err = m.Get(ctx, "key", &rt)
	r.NoError(err)
	r.False(ok)
}
