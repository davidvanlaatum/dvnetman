package logger

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
	"time"
)

func testContext(t testing.TB) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

func TestLogger(t *testing.T) {
	r := require.New(t)
	c := NewCollector()
	l := NewLogger(LevelTrace, c)
	ctx := testContext(t)
	l.Info(ctx).Key("a", "b").Msg("test")
	_, _, line, ok := runtime.Caller(0)
	r.True(ok)
	logs := c.Logs()
	r.NotZero(logs[0].Time)
	r.Equal(
		[]EventData{
			{
				Level:   LevelInfo,
				Message: "test",
				Keys: map[string]interface{}{
					"a": "b",
				},
				driver: l.(*logger),
				Time:   logs[0].Time,
				File:   fmt.Sprintf("pkg/logger/logger_test.go:%d", line-1),
			},
		}, logs,
	)
}

func TestCallsKeyProviders(t *testing.T) {
	r := require.New(t)
	c := NewCollector()
	l := NewLogger(LevelTrace, c)
	ctx := testContext(t)
	called := 0
	ctx = l.SubLogger().KeyProvider(
		func(ctx context.Context) map[string]interface{} {
			called++
			return map[string]interface{}{
				"key": "value",
			}
		},
	).Key("Z", "Z").Logger().Context(ctx)
	Info(ctx).Key("a", "b").Msg("test")
	_, _, line, ok := runtime.Caller(0)
	r.True(ok)
	r.Equal(1, called)
	Ctx(ctx).SubLogger().Key("X", "Y").Logger().Info(ctx).Msg("test2")
	_, _, line2, ok := runtime.Caller(0)
	r.True(ok)
	r.Equal(2, called)
	logs := c.Logs()
	r.Equal(
		EventData{
			Level:   LevelInfo,
			Message: "test",
			Keys: map[string]interface{}{
				"a":   "b",
				"key": "value",
				"Z":   "Z",
			},
			driver: Ctx(ctx).(*logger),
			Time:   logs[0].Time,
			File:   fmt.Sprintf("pkg/logger/logger_test.go:%d", line-1),
		}, logs[0],
	)
	r.Equal(
		EventData{
			Level:   LevelInfo,
			Message: "test2",
			Keys: map[string]interface{}{
				"key": "value",
				"X":   "Y",
				"Z":   "Z",
			},
			driver: logs[1].driver,
			Time:   logs[1].Time,
			File:   fmt.Sprintf("pkg/logger/logger_test.go:%d", line2-1),
		},
		logs[1],
	)
	r.Len(logs, 2)
}

func TestConsoleFormatter_Format(t *testing.T) {
	r := require.New(t)
	f := NewConsoleFormatter()
	testTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	s := f.Format(
		&EventData{
			Time:    testTime,
			Level:   LevelInfo,
			Message: "test",
			Keys: map[string]interface{}{
				"a": "b",
				"c": "d",
			},
			File: "bla.go:123",
		},
	)
	r.Equal("2021-01-01T00:00:00Z bla.go:123 info test a=b c=d", s)
}
