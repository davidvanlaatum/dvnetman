package utils

import (
	"context"
	"github.com/pkg/errors"
)

func PropagateError(f func() error, rtErr *error, msg string) {
	if err := f(); err != nil && *rtErr == nil {
		*rtErr = errors.Wrap(err, msg)
	}
}

func PropagateErrorContext(ctx context.Context, f func(ctx context.Context) error, rtErr *error, msg string) {
	PropagateError(func() error { return f(ctx) }, rtErr, msg)
}
