package cache

import (
	"context"
	"net/url"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, value interface{}) (bool, error)
	Delete(ctx context.Context, key string) error
	Close(ctx context.Context) error
	Flush(ctx context.Context) error
}

type Pool interface {
	New(ctx context.Context) (Cache, error)
	Shutdown(ctx context.Context) error
	Cache
}

type Driver interface {
	New(ctx context.Context, config *url.URL) (Cache, error)
	Name() string
}
