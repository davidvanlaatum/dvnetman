package logger

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	r := require.New(t)
	c := NewCollector()
	l := NewLogger(LevelTrace, c)
	l.Info().Key("a", "b").Msg("test")
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
