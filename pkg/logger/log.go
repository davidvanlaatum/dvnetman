package logger

import (
	"context"
)

type Logger interface {
	Trace(ctx context.Context) Event
	Info(ctx context.Context) Event
	Debug(ctx context.Context) Event
	Warn(ctx context.Context) Event
	Error(ctx context.Context) Event
	Fatal(ctx context.Context) Event
	SubLogger() SubLogger
	Context(ctx context.Context) context.Context
}

type SubLogger interface {
	Key(key string, value interface{}) SubLogger
	KeyProvider(provider func(ctx context.Context) map[string]interface{}) SubLogger
	Logger() Logger
}

type Driver interface {
	Log(data *EventData)
	HelperFuncs() []func()
}

type logger struct {
	drivers      []Driver
	level        Level
	keys         map[string]interface{}
	keyProviders []func(ctx context.Context) map[string]interface{}
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

func (l *logger) getKeys(ctx context.Context) map[string]interface{} {
	keys := make(map[string]interface{}, len(l.keys))
	for k, v := range l.keys {
		keys[k] = v
	}
	for _, provider := range l.keyProviders {
		for k, v := range provider(ctx) {
			keys[k] = v
		}
	}
	return keys
}

func (l *logger) Trace(ctx context.Context) Event {
	if l.level <= LevelTrace {
		return NewEvent(l, LevelTrace, l.getKeys(ctx))
	}
	return &noOpEvent{}
}

func (l *logger) Debug(ctx context.Context) Event {
	if l.level <= LevelDebug {
		return NewEvent(l, LevelDebug, l.getKeys(ctx))
	}
	return &noOpEvent{}
}

func (l *logger) Info(ctx context.Context) Event {
	if l.level <= LevelInfo {
		return NewEvent(l, LevelInfo, l.getKeys(ctx))
	}
	return &noOpEvent{}
}

func (l *logger) Warn(ctx context.Context) Event {
	if l.level <= LevelWarn {
		return NewEvent(l, LevelWarn, l.getKeys(ctx))
	}
	return &noOpEvent{}
}

func (l *logger) Error(ctx context.Context) Event {
	if l.level <= LevelError {
		return NewEvent(l, LevelError, l.getKeys(ctx))
	}
	return &noOpEvent{}
}

func (l *logger) Fatal(ctx context.Context) Event {
	if l.level <= LevelFatal {
		return NewEvent(l, LevelFatal, l.getKeys(ctx))
	}
	return &noOpEvent{}
}

type subLogger struct {
	logger       *logger
	keys         map[string]interface{}
	keyProviders []func(ctx context.Context) map[string]interface{}
}

func (s *subLogger) Key(key string, value interface{}) SubLogger {
	s.keys[key] = value
	return s
}

func (s *subLogger) KeyProvider(provider func(ctx context.Context) map[string]interface{}) SubLogger {
	s.keyProviders = append(s.keyProviders, provider)
	return s
}

func (s *subLogger) Logger() Logger {
	return &logger{
		drivers:      s.logger.drivers,
		level:        s.logger.level,
		keys:         s.keys,
		keyProviders: s.keyProviders,
	}
}

func (l *logger) SubLogger() SubLogger {
	s := &subLogger{
		logger:       l,
		keys:         make(map[string]interface{}, len(l.keys)),
		keyProviders: make([]func(ctx context.Context) map[string]interface{}, len(l.keyProviders)),
	}
	for k, v := range l.keys {
		s.keys[k] = v
	}
	copy(s.keyProviders, l.keyProviders)
	return s
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

func Trace(ctx context.Context) Event {
	return Ctx(ctx).Trace(ctx)
}

func Debug(ctx context.Context) Event {
	return Ctx(ctx).Debug(ctx)
}

func Info(ctx context.Context) Event {
	return Ctx(ctx).Info(ctx)
}

func Warn(ctx context.Context) Event {
	return Ctx(ctx).Warn(ctx)
}

func Error(ctx context.Context) Event {
	return Ctx(ctx).Error(ctx)
}

func Fatal(ctx context.Context) Event {
	return Ctx(ctx).Fatal(ctx)
}
