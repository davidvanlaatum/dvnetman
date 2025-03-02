package cache

import (
	"context"
)

func Lazy[T any](ctx context.Context, c Cache, key string, f func(context.Context) (T, error)) (t T, err error) {
	var ok bool
	if ok, err = c.Get(ctx, key, &t); err != nil {
		return t, err
	} else if ok {
		return t, nil
	}
	if t, err = f(ctx); err != nil {
		return
	} else if err = c.Set(ctx, key, &t); err != nil {
		return
	}
	return
}
