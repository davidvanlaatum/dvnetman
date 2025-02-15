package testutils

import (
	"context"
	"dvnetman/pkg/logger/testlog"
	"testing"
)

func GetTestContext(t testing.TB) context.Context {
	l := testlog.NewTestLog(t)
	ctx, cancel := context.WithCancel(l.Context(context.Background()))
	t.Cleanup(cancel)
	if tt, ok := t.(*testing.T); ok {
		if d, ok := tt.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
			t.Cleanup(cancel)
		}
	}
	return ctx
}
