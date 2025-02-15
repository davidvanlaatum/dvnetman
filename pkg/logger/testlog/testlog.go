package testlog

import (
	"dvnetman/pkg/logger"
	"testing"
)

type TestLog struct {
	t         testing.TB
	formatter logger.Formatter
}

func (t *TestLog) Log(data *logger.EventData) {
	t.t.Helper()
	t.t.Log(t.formatter.Format(data))
}

func (t *TestLog) HelperFuncs() []func() {
	return []func(){t.t.Helper}
}

func NewTestLog(t testing.TB) logger.Logger {
	return logger.NewLogger(
		logger.LevelTrace, &TestLog{
			t:         t,
			formatter: logger.NewConsoleFormatter().DisableCaller(),
		},
	)
}
