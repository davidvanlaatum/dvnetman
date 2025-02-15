package logger

import (
	"fmt"
	"github.com/fatih/color"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var pathPrefix string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		pathPrefix = filepath.Dir(filepath.Dir(filepath.Dir(file))) + "/"
	}
}

type Event interface {
	Msg(msg string)
	Msgf(format string, v ...interface{})
	Key(key string, value interface{}) Event
	Enabled() bool
	Err(err error) Event
	Caller(skip int) Event
}

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l Level) ColorString() *color.Color {
	switch l {
	case LevelTrace:
		return color.New(color.FgCyan)
	case LevelDebug:
		return color.New(color.FgBlue)
	case LevelInfo:
		return color.New(color.FgGreen)
	case LevelWarn:
		return color.New(color.FgYellow)
	case LevelError:
		return color.New(color.FgRed)
	case LevelFatal:
		return color.New(color.FgHiRed)
	default:
		return color.New(color.FgWhite)
	}
}

type EventData struct {
	Time    time.Time
	Level   Level
	Message string
	Keys    map[string]interface{}
	driver  Driver
	File    string
}

func NewEvent(driver Driver, level Level, keys map[string]interface{}) Event {
	k := make(map[string]interface{}, len(keys))
	for key, value := range keys {
		k[key] = value
	}
	return &EventData{
		Level:  level,
		driver: driver,
		Time:   time.Now(),
		Keys:   k,
	}
}

func (e *EventData) Caller(skip int) Event {
	if e.File != "" {
		return e
	}
	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		e.File = fmt.Sprintf("%s:%d", strings.TrimPrefix(file, pathPrefix), line)
	}
	return e
}

func (e *EventData) Err(err error) Event {
	return e.Key("error", err)
}

func (e *EventData) Msg(msg string) {
	e.Message = msg
	for _, f := range e.driver.HelperFuncs() {
		f()
	}
	e.Caller(1)
	e.driver.Log(e)
}

func (e *EventData) Msgf(format string, v ...interface{}) {
	e.Message = fmt.Sprintf(format, v...)
	for _, f := range e.driver.HelperFuncs() {
		f()
	}
	e.Caller(1)
	e.driver.Log(e)
}

func (e *EventData) Key(key string, value interface{}) Event {
	if e.Keys == nil {
		e.Keys = make(map[string]interface{})
	}
	e.Keys[key] = value
	return e
}

func (e *EventData) Enabled() bool {
	return true
}

var _ Event = (*EventData)(nil)

type noOpEvent struct{}

func (n *noOpEvent) Caller(int) Event {
	return n
}

func (n *noOpEvent) Err(error) Event {
	return n
}

func (n *noOpEvent) Msg(string) {
	// noop
}

func (n *noOpEvent) Msgf(string, ...interface{}) {
	// noop
}

func (n *noOpEvent) Key(string, interface{}) Event {
	return n
}

func (n *noOpEvent) Enabled() bool {
	return false
}

var _ Event = (*noOpEvent)(nil)
