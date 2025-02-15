package logger

import (
	"context"
)

type Logger interface {
	Trace() Event
	Info() Event
	Debug() Event
	Warn() Event
	Error() Event
	Fatal() Event
	SubLogger() SubLogger
	Context(ctx context.Context) context.Context
}

type SubLogger interface {
	Key(key string, value interface{}) SubLogger
	Logger() Logger
}

type Driver interface {
	Log(data *EventData)
	HelperFuncs() []func()
}

type logger struct {
	drivers []Driver
	level   Level
	keys    map[string]interface{}
}

func NewLogger(level Level, drivers ...Driver) Logger {
	return &logger{
		drivers: drivers,
		level:   level,
	}
}

func (l *logger) Log(data *EventData) {
	for _, f := range l.HelperFuncs() {
		f()
	}
	for _, driver := range l.drivers {
		driver.Log(data)
	}
}

func (l *logger) HelperFuncs() (rt []func()) {
	for _, driver := range l.drivers {
		rt = append(rt, driver.HelperFuncs()...)
	}
	return rt
}

func (l *logger) Trace() Event {
	if l.level <= LevelTrace {
		return NewEvent(l, LevelTrace, l.keys)
	}
	return &noOpEvent{}
}

func (l *logger) Debug() Event {
	if l.level <= LevelDebug {
		return NewEvent(l, LevelDebug, l.keys)
	}
	return &noOpEvent{}
}

func (l *logger) Info() Event {
	if l.level <= LevelInfo {
		return NewEvent(l, LevelInfo, l.keys)
	}
	return &noOpEvent{}
}

func (l *logger) Warn() Event {
	if l.level <= LevelWarn {
		return NewEvent(l, LevelWarn, l.keys)
	}
	return &noOpEvent{}
}

func (l *logger) Error() Event {
	if l.level <= LevelError {
		return NewEvent(l, LevelError, l.keys)
	}
	return &noOpEvent{}
}

func (l *logger) Fatal() Event {
	if l.level <= LevelFatal {
		return NewEvent(l, LevelFatal, l.keys)
	}
	return &noOpEvent{}
}

type subLogger struct {
	logger *logger
	keys   map[string]interface{}
}

func (s *subLogger) Key(key string, value interface{}) SubLogger {
	s.keys[key] = value
	return s
}

func (s *subLogger) Logger() Logger {
	return &logger{
		drivers: s.logger.drivers,
		level:   s.logger.level,
		keys:    s.keys,
	}
}

func (l *logger) SubLogger() SubLogger {
	return &subLogger{
		logger: l,
		keys:   map[string]interface{}{},
	}
}

type contextKeyType struct{}

var contextKey = &contextKeyType{}

func (l *logger) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, l)
}

var _ Logger = (*logger)(nil)
var _ Driver = (*logger)(nil)

func Ctx(ctx context.Context) Logger {
	l := ctx.Value(contextKey)
	if x, ok := l.(Logger); ok {
		return x
	}
	panic("logger not found in context")
}

//func LogFromContext(ctx context.Context, logger Logger) Logger {
//	span := trace.SpanFromContext(ctx)
//	if !span.SpanContext().IsValid() {
//		return logger
//	}
//	prefix := span.SpanContext().TraceID().String() + "/" + span.SpanContext().SpanID().String()
//	if sl, ok := logger.(SubLogger); ok {
//		return sl.GetSubLogger(prefix)
//	}
//	if sl, ok := logger.(*log.Logger); ok {
//		return log.New(sl.Writer(), prefix+" ", sl.Flags())
//	}
//	return logger
//}
